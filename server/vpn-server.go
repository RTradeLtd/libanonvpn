package samforwardervpnserver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	//	"time"

	whitelister "github.com/eyedeekay/accessregister/auth"
	//"github.com/eyedeekay/canal"
	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	sfi2pkeys "github.com/eyedeekay/sam-forwarder/i2pkeys"
	samforwarder "github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
	udptunnel "github.com/eyedeekay/udptunnel/tunnel"
)

// SAMClientServerVPN is a VPN service which provides access over UDP-like
// I2P connections to provide clients with anonymity, authentication, and
// encryption
type SAMClientServerVPN struct {
	// VPN tunnel
	*udptunnel.Tunnel
	//*wallet.InnerWallet
	WhiteListers []whitelister.WhiteLister
	*hashhash.Hasher
	*i2ptunconf.Conf
	whiteListingTunnel *samforwarder.SAMForwarder

	// i2p tunnel configuration
	samConn *sam3.SAM
	ports   []uint16
	keys    i2pkeys.I2PKeys
	up      bool

	// Client Tunnel Configuration
	ClientFilePath string
	eepProxyString string
}

// Check implements a WhiteLister that runs the internal whitelisters
func (w *SAMClientServerVPN) Check(x interface{}) (string, interface{}, bool) {
	if len(w.WhiteListers) == 0 {
		return "", nil, true
	}
	for _, v := range w.WhiteListers {
		base64, extra, outcome := v.Check(x)
		if outcome {
			return base64, extra, outcome
		}
	}
	return "", nil, false
}

func (f *SAMClientServerVPN) Config() *i2ptunconf.Conf {
	return f.Conf
}

func (f *SAMClientServerVPN) accesslisttype() string {
	if f.Config().AccessListType == "whitelist" {
		return "i2cp.enableAccessList=true"
	} else if f.Config().AccessListType == "blacklist" {
		return "i2cp.enableBlackList=true"
	} else if f.Config().AccessListType == "none" {
		return ""
	}
	return ""
}

func (f *SAMClientServerVPN) accesslist() string {
	if f.Config().AccessListType != "" && len(f.Config().AccessList) > 0 {
		r := ""
		for _, s := range f.Config().AccessList {
			r += s + ","
		}
		return "i2cp.accessList=" + strings.TrimSuffix(r, ",")
	}
	return ""
}

func (f *SAMClientServerVPN) leasesetsettings() (string, string, string) {
	var r, s, t string
	if f.Config().LeaseSetKey != "" {
		r = "i2cp.leaseSetKey=" + f.Config().LeaseSetKey
	}
	if f.Config().LeaseSetPrivateKey != "" {
		s = "i2cp.leaseSetPrivateKey=" + f.Config().LeaseSetPrivateKey
	}
	if f.Config().LeaseSetPrivateSigningKey != "" {
		t = "i2cp.leaseSetPrivateSigningKey=" + f.Config().LeaseSetPrivateSigningKey
	}
	return r, s, t
}

func (f *SAMClientServerVPN) print() []string {
	lsk, lspk, lspsk := f.leasesetsettings()
	return []string{
		"inbound.length=" + fmt.Sprintf("%d", f.Config().InLength),
		"outbound.length=" + fmt.Sprintf("%d", f.Config().OutLength),
		"inbound.lengthVariance=" + fmt.Sprintf("%d", f.Config().InVariance),
		"outbound.lengthVariance=" + fmt.Sprintf("%d", f.Config().OutVariance),
		"inbound.backupQuantity=" + fmt.Sprintf("%d", f.Config().InBackupQuantity),
		"outbound.backupQuantity=" + fmt.Sprintf("%d", f.Config().OutBackupQuantity),
		"inbound.quantity=" + fmt.Sprintf("%d", f.Config().InQuantity),
		"outbound.quantity=" + fmt.Sprintf("%d", f.Config().OutQuantity),
		"inbound.allowZeroHop=" + fmt.Sprintf("%t", f.Config().InAllowZeroHop),
		"outbound.allowZeroHop=" + fmt.Sprintf("%t", f.Config().OutAllowZeroHop),
		"i2cp.fastRecieve=" + fmt.Sprintf("%t", f.Config().FastRecieve),
		"i2cp.gzip=" + fmt.Sprintf("%t", f.Config().UseCompression),
		"i2cp.reduceOnIdle=" + fmt.Sprintf("%t", f.Config().ReduceIdle),
		"i2cp.reduceIdleTime=" + fmt.Sprintf("%d", f.Config().ReduceIdleTime),
		"i2cp.reduceQuantity=" + fmt.Sprintf("%d", f.Config().ReduceIdleQuantity),
		"i2cp.closeOnIdle=" + fmt.Sprintf("%t", f.Config().CloseIdle),
		"i2cp.closeIdleTime=" + fmt.Sprintf("%d", f.Config().CloseIdleTime),
		"i2cp.messageReliability=" + fmt.Sprintf("%s", f.Config().MessageReliability),
		"i2cp.encryptLeaseSet=" + fmt.Sprintf("%t", f.Config().EncryptLeaseSet),
		lsk, lspk, lspsk,
		f.accesslisttype(),
		f.accesslist(),
	}
}

// Props returns a map of the options in use as strings
func (f *SAMClientServerVPN) Props() map[string]string {
	r := make(map[string]string)
	print := f.print()
	print = append(print, "base32="+f.Base32())
	print = append(print, "base64="+f.Base64())
	print = append(print, "base32words="+f.Base32Readable())
	for _, prop := range print {
		k, v := sfi2pkeys.Prop(prop)
		r[k] = v
	}
	for key, prop := range f.whiteListingTunnel.Props() {
		r["whitelister"+key] = prop
	}
	return r
}

// Print returns the set of options in use by the tunnel as one long string.
func (f *SAMClientServerVPN) Print() string {
	var r string
	r += "name=" + f.Config().TunName + "\n"
	r += "type=" + f.Config().Type + "\n"
	r += "base32=" + f.Base32() + "\n"
	r += "base64=" + f.Base64() + "\n"
	r += "vpnserver\n"
	for _, s := range f.print() {
		r += s + "\n"
	}
	return strings.Replace(r, "\n\n", "\n", -1)
}

// Search searches the tunnel options for a specific property
func (f *SAMClientServerVPN) Search(search string) string {
	terms := strings.Split(search, ",")
	if search == "" {
		return f.Print()
	}
	for _, value := range terms {
		if !strings.Contains(f.Print(), value) {
			return ""
		}
	}
	return f.Print()
}

// Target returns the host:port of the local service you want to forward to i2p
func (f *SAMClientServerVPN) Target() string {
	return f.Config().TunnelHost
}

// SAMSetupSock creates a socket which can bs used by the UDP tunnel
func (f *SAMClientServerVPN) SAMSetupSock(netAddr net.Addr) net.PacketConn {
	sock, err := f.samConn.NewDatagramSession(f.Config().TunName, f.Keys(), f.Config().PrintSlice(), 0)
	if err != nil {
		log.Fatalf("error listening on socket: %v", err)
	}
	err = sock.SetWriteBuffer(14 * 1024)
	if err != nil {
		log.Fatalf("error listening on socket: %v", err)
	} //
	return sock
}

func (f *SAMClientServerVPN) sam() string {
	return f.Config().SamHost + ":" + f.Config().SamPort
}

// TargetAddr returns the sam3.I2PAddr of the VPN Service tunnel as a net.Addr
func (f *SAMClientServerVPN) TargetAddr() (net.Addr, error) {
	return f.Keys().Addr(), nil
}

// LocalAddr returns the IP address of the VPN Service tunnel, a.k.a. the
// Gateway. This is *not* the public IP address of the VPN Service
func (f *SAMClientServerVPN) LocalAddr() (net.Addr, error) {
	log.Println(f.Config().TunnelHost)
	return &net.IPAddr{IP: net.ParseIP(f.Config().TunnelHost)}, nil
}

// Base32 returns the base32 address of the local tunnel, *not* the destination
// tunnel
func (f *SAMClientServerVPN) Base32() string {
	return f.Keys().Addr().Base32()
}

//Base32Readable returns the base32 address where the local service is being forwarded
func (f *SAMClientServerVPN) Base32Readable() string {
	b32 := strings.Replace(f.Base32(), ".b32.i2p", "", 1)
	rhash, _ := f.Hasher.Friendly(b32)
	return rhash + " " + strconv.Itoa(len(b32))
}

// Base64 returns the base64 address of the local tunnel, *not* the destination
// tunnel
func (f *SAMClientServerVPN) Base64() string {
	return f.Keys().Addr().Base64()
}

// Cleanup does nothing in this case, because AFAIK you should not use these in
// half-loaded state ever and this everything that would normally be in Cleanup
// is in Close instead, which normally calls Cleanup before doing a closing
// operation.
func (f *SAMClientServerVPN) Cleanup() {
	return
}

// Close closes the underlying SAM Connection to I2P
func (f *SAMClientServerVPN) Close() error {
	return f.samConn.Close()
}

//
func (f *SAMClientServerVPN) GetType() string {
	return "vpnserver"
}

func (f *SAMClientServerVPN) ID() string {
	return f.Config().TunName
}

func (f *SAMClientServerVPN) Keys() i2pkeys.I2PKeys {
	return f.keys
}

func (f *SAMClientServerVPN) Up() bool {
	return f.up
}

func (f *SAMClientServerVPN) ClientConfig() string {
	label := "[" + f.ID() + "]\n"
	label += "type = vpnclient\n"
	label += "destination = " + f.Base32() + "\n"
	return label
}

func (f *SAMClientServerVPN) LoadRemoteAddr() net.Addr {
	addr := f.Tunnel.RemoteAddr.Load()
	switch f.Tunnel.RemoteAddr.Load().(type) {
	case nil:
		return nil
	case i2pkeys.I2PAddr:
		k := addr.(i2pkeys.I2PAddr)
		return &k
	case *i2pkeys.I2PAddr:
		return addr.(i2pkeys.I2PAddr)
	}
	return nil
}

func (f *SAMClientServerVPN) UpdateRemoteAddr(addr net.Addr) {
	oldAddr := f.LoadRemoteAddr()
	log.Println("Finding new endpoint")
	if addr == nil {
		return
	}
	if oldAddr == nil {
		f.Tunnel.RemoteAddr.Store(addr)
		log.Printf("switching remote address: %v != %v", addr, oldAddr)
		return
	}
	if addr != nil && (oldAddr == nil || addr.String() != oldAddr.String()) { //|| addr.Zone != oldAddr.Zone) {
		f.Tunnel.RemoteAddr.Store(addr)
		log.Printf("switching remote address: %v != %v", addr, oldAddr)
	}
}

func (t *SAMClientServerVPN) WriteConn(sock net.PacketConn, raddr net.Addr, b []byte, n int, magic [16]byte) (int, error) {
	if raddr.(*i2pkeys.I2PAddr) == nil {
		return 0, nil
	}
	if n != 0 {
		return sock.(*sam3.DatagramSession).WriteTo(b[:n], raddr.(*i2pkeys.I2PAddr))
	}
	return sock.(*sam3.DatagramSession).WriteTo(magic[:], raddr.(*i2pkeys.I2PAddr))
}

func (s *SAMClientServerVPN) Serve() error {
	var logBuf bytes.Buffer
	Logger := log.New(io.MultiWriter(os.Stderr, &logBuf), "", log.Ldate|log.Ltime|log.Lshortfile)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		Logger.Printf("received %v - initiating shutdown", <-sigc)
		cancel()
	}()
	localaddr, err := s.LocalAddr()
	if err != nil {
		return err
	}
	targetaddr, err := s.TargetAddr()
	if err != nil {
		return err
	}
	log.Println("My base64 public key is:", targetaddr.(i2pkeys.I2PAddr).Base64())
	if len(s.ports) == 0 {
		s.ports = append(s.ports, 22)
	}
	s.Tunnel = udptunnel.NewCustomTunnel(
		true,               // is a server
		s.Config().TunName, // device is shared with i2p tunnel name
		localaddr,          // tunLocalAddr
		targetaddr,         // tunRemoteAddr
		targetaddr,         // netAddr
		s.ports,            // ports allowed
		"i2p",              // magic word
		30,                 // heartbeat interval
		Logger,             // logger
		s.SAMSetupSock,     // Socket Configurator
		s.TargetAddr,       //
		s.UpdateRemoteAddr, //
		s.LoadRemoteAddr,   //
		s.WriteConn,        //
	)
	if err := s.Tunnel.Setup(); err != nil {
		return err
	}
	go s.whiteListingTunnel.Serve()
	/*if err = firewall.ServerSetup(s.Config().TunName, "eth0"); err != nil {
		return err
	}*/
	s.Tunnel.Run(ctx)

	return nil
}
