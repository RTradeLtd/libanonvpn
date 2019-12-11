package samforwardervpn

import (
	"fmt"
	"log"
	"strings"

	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	sfi2pkeys "github.com/eyedeekay/sam-forwarder/i2pkeys"
	samtunnel "github.com/eyedeekay/sam-forwarder/interface"
	"github.com/eyedeekay/sam3"
)

func (s *SAMClientVPN) Load() (samtunnel.SAMTunnel, error) {
	var err error
	s.samConn, err = sam3.NewSAM(s.sam())
	if err != nil {
		return nil, err
	}
	log.Println("SAM Bridge connection established.")
	if s.Conf.SaveFile {
		log.Println("Saving i2p keys")
	}
	if s.keys, err = sfi2pkeys.Load(s.Conf.FilePath, s.Conf.TunName, s.Conf.KeyFilePath, s.samConn, s.Conf.SaveFile); err != nil {
		return nil, err
	}
	log.Println("Destination keys generated, tunnel name:", s.Conf.TunName)
	if s.Conf.SaveFile {
		if err := sfi2pkeys.Save(s.Conf.FilePath, s.Conf.TunName, s.Conf.KeyFilePath, s.Keys()); err != nil {
			return nil, err
		}
		log.Println("Saved tunnel keys for", s.Conf.TunName)
	}
	s.up = true
	/*err = s.InnerWallet.LoadWallet()*/
	if err != nil {
		return nil, err
	}
	s.Hasher, err = hashhash.NewHasher(len(strings.Replace(s.Base32(), ".b32.i2p", "", 1)))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewSAMClientVPN(conf *i2ptunconf.Conf, destination ...string) (*SAMClientVPN, error) {
	if len(destination) == 0 {
		return NewSAMClientVPNFromOptions(SetClientVPNConfig(conf))
	} else if len(destination) == 1 {
		return NewSAMClientVPNFromOptions(SetClientVPNConfig(conf), SetClientDest(destination[0]))
	} else {
		return nil, fmt.Errorf("Error, argument for destination must be len==0 or len==1")
	}
}

func NewSAMClientVPNFromOptions(opts ...func(*SAMClientVPN) error) (*SAMClientVPN, error) {
	var s SAMClientVPN
	s.Conf = &i2ptunconf.Conf{}
	s.FilePath = ""
	//s.InnerWallet = &innerwallet.InnerWallet{}
	//s.InnerWallet.MultiWalletConf = &multiwalletconf.Config{}
	//s.InnerWallet.ChainCfg = &chaincfg.MainNetParams
	for _, o := range opts {
		if err := o(&s); err != nil {
			return nil, err
		}
	}
	var err error
	if s.Conf == nil {
		if s.FilePath != "" {
			s.Conf, err = i2ptunconf.NewI2PTunConf(s.FilePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("No VPN configuration provided")
		}
	}
	r, e := s.Load()
	if e != nil {
		return nil, e
	}
	return r.(*SAMClientVPN), nil
}

// NewSAMVPNForwarderFromConfig generates a new SAMVPNForwarder from a config file
func NewSAMVPNClientForwarderFromConfig(iniFile, SamHost, SamPort string, label ...string) (*SAMClientVPN, error) {
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
		return NewSAMClientVPN(config)
	}
	return nil, nil
}
