package samforwardervpnserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	whitelister "github.com/eyedeekay/accessregister/auth"
	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	sfi2pkeys "github.com/eyedeekay/sam-forwarder/i2pkeys"
	samtunnel "github.com/eyedeekay/sam-forwarder/interface"
	samforwarder "github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam3"
	"github.com/phayes/freeport"
)

func (s *SAMClientServerVPN) Load() (samtunnel.SAMTunnel, error) {
	var err error
	//SAM Initialization
	s.samConn, err = sam3.NewSAM(s.sam())
	if err != nil {
		return nil, err
	}
	log.Println("SAM Bridge connection established.")
	if s.Config().SaveFile {
		log.Println("Saving i2p keys")
	}
	if s.keys, err = sfi2pkeys.Load(s.Config().SaveDirectory, s.Config().TunName, s.Config().FilePath, s.samConn, s.Config().SaveFile); err != nil {
		return nil, err
	}
	log.Println("Destination keys generated, tunnel name:", s.Config().TunName)
	if s.Config().SaveFile {
		if err := sfi2pkeys.Save("", s.Config().TunName, s.Config().KeyFilePath, s.Keys()); err != nil {
			return nil, err
		}
		log.Println("Saved tunnel keys for", s.Config().TunName)
	}
	s.up = true
	if s.ClientFilePath == "" {
		s.ClientFilePath = s.Config().TunName + ".ini"
	}
	err = ioutil.WriteFile(s.ClientFilePath, []byte(s.ClientConfig()), 0644)
	if err != nil {
		return nil, err
	}
	/*tmp, err := innerwallet.NewWallet()
	if err != nil {
		return nil, err
	}
	s.WhiteListers = append(s.WhiteListers, tmp)
	err = s.InnerWallet.LoadWallet()
	if err != nil {
		return nil, err
	}*/
	s.Hasher, err = hashhash.NewHasher(len(strings.Replace(s.Base32(), ".b32.i2p", "", 1)))
	if err != nil {
		return nil, err
	}
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	s.whiteListingTunnel, err = samforwarder.NewSAMForwarderFromOptions(
		samforwarder.SetType("httpserver"),
		samforwarder.SetSaveFile(s.Config().SaveFile),
		samforwarder.SetFilePath(s.Config().SaveDirectory),
		samforwarder.SetHost("localhost"),
		samforwarder.SetPort(strconv.Itoa(port)),
		samforwarder.SetSAMHost(s.Config().SamHost),
		samforwarder.SetSAMPort(s.Config().SamPort),
		samforwarder.SetSigType(s.Config().SigType),
		samforwarder.SetName("whitelister."+s.Config().TunName),
		samforwarder.SetInLength(s.Config().InLength),
		samforwarder.SetOutLength(s.Config().OutLength),
		samforwarder.SetInVariance(s.Config().InVariance),
		samforwarder.SetOutVariance(s.Config().OutVariance),
		samforwarder.SetInQuantity(s.Config().InQuantity),
		samforwarder.SetOutQuantity(s.Config().OutQuantity),
		samforwarder.SetInBackups(s.Config().InBackupQuantity),
		samforwarder.SetOutBackups(s.Config().OutBackupQuantity),
		samforwarder.SetEncrypt(s.Config().EncryptLeaseSet),
		samforwarder.SetLeaseSetKey(s.Config().LeaseSetKey),
		samforwarder.SetLeaseSetPrivateKey(s.Config().LeaseSetPrivateKey),
		samforwarder.SetLeaseSetPrivateSigningKey(s.Config().LeaseSetPrivateSigningKey),
		samforwarder.SetAllowZeroIn(s.Config().InAllowZeroHop),
		samforwarder.SetAllowZeroOut(s.Config().OutAllowZeroHop),
		samforwarder.SetFastRecieve(s.Config().FastRecieve),
		samforwarder.SetCompress(s.Config().UseCompression),
		samforwarder.SetReduceIdle(s.Config().ReduceIdle),
		samforwarder.SetReduceIdleTimeMs(s.Config().ReduceIdleTime),
		samforwarder.SetReduceIdleQuantity(s.Config().ReduceIdleQuantity),
		samforwarder.SetCloseIdle(s.Config().CloseIdle),
		samforwarder.SetCloseIdleTimeMs(s.Config().CloseIdleTime),
		samforwarder.SetAccessListType(s.Config().AccessListType),
		samforwarder.SetAccessList(s.Config().AccessList),
		samforwarder.SetMessageReliability(s.Config().MessageReliability),
		samforwarder.SetKeyFile(s.Config().KeyFilePath),
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewSAMClientServerVPNFromOptions(opts ...func(*SAMClientServerVPN) error) (*SAMClientServerVPN, error) {
	var s SAMClientServerVPN
	s.Conf = &i2ptunconf.Conf{}
	//s.InnerWallet = &innerwallet.InnerWallet{}
	//s.InnerWallet.MultiWalletConf = &multiwalletconf.Config{}
	//s.InnerWallet.ChainCfg = &chaincfg.MainNetParams
	s.WhiteListers = []whitelister.WhiteLister{}
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	var err error
	if s.Config().FilePath != "" {
		s.Conf, err = i2ptunconf.NewI2PTunConf(s.Config().FilePath)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No VPN configuration provided")
	}
	r, e := s.Load()
	if e != nil {
		return nil, e
	}
	return r.(*SAMClientServerVPN), nil
}

// NewSAMVPNClientForwarderFromConfig generates a new SAMVPNForwarder from a config file
func NewSAMVPNForwarderFromConfig(iniFile, SamHost, SamPort string, label ...string) (*SAMClientServerVPN, error) {
	if iniFile != "none" {
		config, err := i2ptunconf.NewI2PTunConf(iniFile, label...)
		if err != nil {
			return nil, err
		}
		if SamHost != "" && SamHost != "127.0.0.1" && SamHost != "localhost" {
			config.SamHost = config.GetSAMHost(SamHost, config.SamHost)
		}
		if SamPort != "" && SamPort != "7656" {
			config.SamPort = config.GetSAMPort(SamPort, config.SamPort)
		}
		return NewSAMClientServerVPN(config)
	}
	return nil, nil
}

func NewSAMClientServerVPN(conf *i2ptunconf.Conf) (*SAMClientServerVPN, error) {
	return NewSAMClientServerVPNFromOptions(SetVPNConfig(conf))
}
