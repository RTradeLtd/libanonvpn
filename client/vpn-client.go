package samforwardervpn

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

	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	sfi2pkeys "github.com/eyedeekay/sam-forwarder/i2pkeys"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
	udptunnel "github.com/eyedeekay/udptunnel/tunnel"
)

// SAMClientVPN is a VPN client which uses UDP-like communication over I2P for
// anonymity, authentication and encryption.
type SAMClientVPN struct {
	// VPN tunnel
	*udptunnel.Tunnel
	//*wallet.InnerWallet
	*hashhash.Hasher
	socksTunnel *samforwarder.SAMForwarder

	// i2p tunnel
	samConn        *sam3.SAM
	Conf           *i2ptunconf.Conf
	FilePath       string
	ports          []uint16
	keys           i2pkeys.I2PKeys
	up             bool
	eepProxyString string
}

func (f *SAMClientVPN) Config() *i2ptunconf.Conf {
	return f.Conf
}

func (f *SAMClientVPN) accesslisttype() string {
	if f.Config().AccessListType == "whitelist" {
		return "i2cp.enableAccessList=true"
	} else if f.Config().AccessListType == "blacklist" {
		return "i2cp.enableBlackList=true"
	} else if f.Config().AccessListType == "none" {
		return ""
	}
	return ""
}

func (f *SAMClientVPN) accesslist() string {
	if f.Config().AccessListType != "" && len(f.Config().AccessList) > 0 {
		r := ""
		for _, s := range f.Config().AccessList {
			r += s + ","
		}
		return "i2cp.accessList=" + strings.TrimSuffix(r, ",")
	}
	return ""
}

func (f *SAMClientVPN) leasesetsettings() (string, string, string) {
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

func (f *SAMClientVPN) print() []string {
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
func (f *SAMClientVPN) Props() map[string]string {
	r := make(map[string]string)
	print := f.print()
	print = append(print, "base32="+f.Base32())
	print = append(print, "base64="+f.Base64())
	print = append(print, "base32words="+f.Base32Readable())
	for _, prop := range print {
		k, v := sfi2pkeys.Prop(prop)
		r[k] = v
	}
	return r
}

// Print returns the set of options in use by the tunnel as one long string.
func (f *SAMClientVPN) Print() string {
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
func (f *SAMClientVPN) Search(search string) string {
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

func (f *SAMClientVPN) AddAccessListMemberAddr(addr i2pkeys.I2PAddr) error {
	return nil
}

func (f *SAMClientVPN) RequestAdditionToAccessList() error {
	return nil
}

func (f *SAMClientVPN) sam() string {
	return f.Config().SamHost + ":" + f.Config().SamPort
}

// Target returns the host IP of the VPN Client
func (f *SAMClientVPN) Target() string {
	addr, _ := f.LocalAddr()
	return addr.String()
}

// SAMSetupSock creates a socket which can bs used by the UDP tunnel
func (f *SAMClientVPN) SAMSetupSock(netAddr net.Addr) net.PacketConn {
	log.Println("Setting up socket connection to", f.Keys())
	sock, err := f.samConn.NewDatagramSession(f.Config().TunName, f.Keys(), f.Config().PrintSlice(), 0)
	if err != nil {
		log.Fatalf("error listening on socket: %v", err)
	}
	err = sock.SetWriteBuffer(14 * 1024)
	if err != nil {
		log.Fatalf("error listening on socket: %v", err)
	}
	log.Println("Listening on SAM Socket", f.Config().TunName)
	return sock
}

// DestinationAddr does the full lookup of the client destination and returns
// the corresponding sam3.I2PAddr as a net.Addr
func (f *SAMClientVPN) DestinationAddr() (net.Addr, error) {
	//log.Println("Looking up gateway destination Base64 for", f.Config().ClientDest)
	addr, err := f.samConn.Lookup(f.Config().ClientDest)
	if err != nil {
		return nil, err
	}
	//log.Println("Gateway VPN Tunnel Address", addr.String())
	return &addr, nil
}

// LocalAddr converts the tunnel IP address string to a net.IPAddr and returns
// it as a net.Addr
func (f *SAMClientVPN) LocalAddr() (net.Addr, error) {
	if f.Config().TargetHost == "127.0.0.1" {
		f.Config().TargetHost = "10.79.0.2"
	}
	log.Println("Local VPN Tunnel Host:", f.Config().TargetHost)
	return &net.IPAddr{IP: net.ParseIP(f.Config().TargetHost)}, nil
}

// RemoteAddr converts the tunnel IP address string to a net.IPAddr and returns
// it as a net.Addr
func (f *SAMClientVPN) RemoteAddr() (net.Addr, error) {
	log.Println("Remote VPN Tunnel Host:", f.Config().TunnelHost)
	return &net.IPAddr{IP: net.ParseIP(f.Config().TunnelHost)}, nil //net.ResolveIPAddr("ip", f.Config().TunnelHost)
}

// Base32 returns the base32 address of the local tunnel, *not* the destination
// tunnel
func (f *SAMClientVPN) Base32() string {
	return f.Keys().Addr().Base32()
}

//Base32Readable returns the base32 address where the local service is being forwarded
func (f *SAMClientVPN) Base32Readable() string {
	b32 := strings.Replace(f.Base32(), ".b32.i2p", "", 1)
	rhash, _ := f.Hasher.Friendly(b32)
	return rhash + " " + strconv.Itoa(len(b32))
}

// Base64 returns the base64 address of the local tunnel, *not* the destination
// tunnel
func (f *SAMClientVPN) Base64() string {
	return f.Keys().Addr().Base64()
}

// Cleanup does nothing in this case, because AFAIK you should not use these in
// half-loaded state ever and this everything that would normally be in Cleanup
// is in Close instead, which normally calls Cleanup before doing a closing
// operation.
func (f *SAMClientVPN) Cleanup() {
	return
}

func (f *SAMClientVPN) Close() error {
	return f.samConn.Close()
}

func (f *SAMClientVPN) GetType() string {
	return "vpnclient"
}

func (f *SAMClientVPN) Keys() i2pkeys.I2PKeys {
	return f.keys
}

func (f *SAMClientVPN) ID() string {
	return f.Config().TunName
}

func (f *SAMClientVPN) Up() bool {
	return f.up
}

func (f *SAMClientVPN) LoadRemoteAddr() net.Addr {
	addr, _ := f.DestinationAddr()
	return addr
}

func (f *SAMClientVPN) UpdateRemoteAddr(addr net.Addr) {
	oldAddr := f.LoadRemoteAddr()
	if addr != nil && (oldAddr == nil || addr.String() != oldAddr.String()) { //|| addr.Zone != oldAddr.Zone) {
		f.Tunnel.RemoteAddr.Store(addr)
		log.Printf("switching remote address: %v != %v", addr, oldAddr)
	}
}

func (t *SAMClientVPN) WriteConn(sock net.PacketConn, raddr net.Addr, b []byte, n int, magic [16]byte) (int, error) {
	if raddr.(*i2pkeys.I2PAddr) == nil {
		return 0, nil
	}
	if n != 0 {
		return sock.(*sam3.DatagramSession).WriteTo(b[:n], raddr.(*i2pkeys.I2PAddr))
	}
	return sock.(*sam3.DatagramSession).WriteTo(magic[:], raddr.(*i2pkeys.I2PAddr))
}

func (s *SAMClientVPN) Serve() error {
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
	log.Println("local addr:", localaddr)
	destinationaddr, err := s.DestinationAddr()
	if err != nil {
		return err
	}
	log.Println("gateway destination:", destinationaddr)
	remoteaddr, err := s.RemoteAddr()
	if err != nil {
		return err
	}
	log.Println("gateway addr:", remoteaddr)
	if len(s.ports) == 0 {
		s.ports = append(s.ports, 22)
	}
	log.Println("Setting up VPN Tunnel")
	s.Tunnel = udptunnel.NewCustomTunnel(
		false,              // is a server
		s.Config().TunName, // device is shared with i2p tunnel name
		localaddr,          // tunLocalAddr
		remoteaddr,         // tunRemoteAddr
		destinationaddr,    // netAddr
		s.ports,            // ports allowed
		"i2p",              // magic word
		30,                 // heartbeat interval
		Logger,             // logger
		s.SAMSetupSock,     // Socket Configurator
		s.DestinationAddr,  //
		s.UpdateRemoteAddr, //
		s.LoadRemoteAddr,   //
		s.WriteConn,        //
	)
	log.Println("firewall: Setting up Interfaces")
	if err := s.Tunnel.Setup(); err != nil {
		return err
	}
	s.Tunnel.Run(ctx)
	for {

	}
	return nil
}
