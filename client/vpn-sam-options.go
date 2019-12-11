package samforwardervpn

import (
	"fmt"
	"strconv"

	i2ptunconf "github.com/eyedeekay/sam-forwarder/config"
)

//ClientOption is a SAMClientVPN Option
type ClientOption func(*SAMClientVPN) error

func SetClientFilePath(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.FilePath = s
		return nil
	}
}

func SetEepProxy(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.eepProxyString = s
		return nil
	}
}

func SetClientVPNConfig(s *i2ptunconf.Conf) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf = s
		return nil
	}
}

func SetClientDest(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.ClientDest = s
		return nil
	}
}

//SetSigType sets the type of the forwarder server
func SetSigType(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if s == "" {
			c.Conf.SigType = ""
		} else if s == "DSA_SHA1" {
			c.Conf.SigType = "DSA_SHA1"
		} else if s == "ECDSA_SHA256_P256" {
			c.Conf.SigType = "ECDSA_SHA256_P256"
		} else if s == "ECDSA_SHA384_P384" {
			c.Conf.SigType = "ECDSA_SHA384_P384"
		} else if s == "ECDSA_SHA512_P521" {
			c.Conf.SigType = "ECDSA_SHA512_P521"
		} else if s == "EdDSA_SHA512_Ed25519" {
			c.Conf.SigType = "EdDSA_SHA512_Ed25519"
		} else {
			c.Conf.SigType = "EdDSA_SHA512_Ed25519"
		}
		return nil
	}
}

//SetType sets the type of the forwarder server
func SetType(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if s == "client" {
			c.Conf.Type = s
			return nil
		} else {
			c.Conf.Type = "server"
			return nil
		}
	}
}

//SetSaveFile tells the router to save the tunnel's keys long-term
func SetSaveFile(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.SaveFile = b
		return nil
	}
}

//SetPointHost sets the VPN host of the local machine
func SetPointHost(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.TargetHost = s
		return nil
	}
}

//SetEndpointHost sets the host of the service to forward
func SetEndpointHost(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.TunnelHost = s
		return nil
	}
}

//SetSAMHost sets the host of the SAMClientVPN's SAM bridge
func SetSAMHost(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.SamHost = s
		return nil
	}
}

//SetSAMPort sets the port of the SAMClientVPN's SAM bridge using a string
func SetSAMPort(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid SAM Port %s; non-number", s)
		}
		if port < 65536 && port > -1 {
			c.Conf.SamPort = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetName sets the host of the SAMClientVPN's SAM bridge
func SetName(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.TunName = s
		return nil
	}
}

//SetInLength sets the number of hops inbound
func SetInLength(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u < 7 && u >= 0 {
			c.Conf.InLength = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutLength sets the number of hops outbound
func SetOutLength(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u < 7 && u >= 0 {
			c.Conf.OutLength = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel length")
	}
}

//SetInVariance sets the variance of a number of hops inbound
func SetInVariance(i int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if i < 7 && i > -7 {
			c.Conf.InVariance = i
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutVariance sets the variance of a number of hops outbound
func SetOutVariance(i int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if i < 7 && i > -7 {
			c.Conf.OutVariance = i
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel variance")
	}
}

//SetInQuantity sets the inbound tunnel quantity
func SetInQuantity(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u <= 16 && u > 0 {
			c.Conf.InQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel quantity")
	}
}

//SetOutQuantity sets the outbound tunnel quantity
func SetOutQuantity(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u <= 16 && u > 0 {
			c.Conf.OutQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel quantity")
	}
}

//SetInBackups sets the inbound tunnel backups
func SetInBackups(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u < 6 && u >= 0 {
			c.Conf.InBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel backup quantity")
	}
}

//SetOutBackups sets the inbound tunnel backups
func SetOutBackups(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u < 6 && u >= 0 {
			c.Conf.OutBackupQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel backup quantity")
	}
}

//SetEncrypt tells the router to use an encrypted leaseset
func SetEncrypt(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.EncryptLeaseSet = b
		return nil
	}
}

//SetLeaseSetKey sets the host of the SAMClientVPN's SAM bridge
func SetLeaseSetKey(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.LeaseSetKey = s
		return nil
	}
}

//SetLeaseSetPrivateKey sets the host of the SAMClientVPN's SAM bridge
func SetLeaseSetPrivateKey(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.LeaseSetPrivateKey = s
		return nil
	}
}

//SetLeaseSetPrivateSigningKey sets the host of the SAMClientVPN's SAM bridge
func SetLeaseSetPrivateSigningKey(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.LeaseSetPrivateSigningKey = s
		return nil
	}
}

//SetMessageReliability sets the host of the SAMClientVPN's SAM bridge
func SetMessageReliability(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.MessageReliability = s
		return nil
	}
}

//SetAllowZeroIn tells the tunnel to accept zero-hop peers
func SetAllowZeroIn(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.InAllowZeroHop = b
		return nil
	}
}

//SetAllowZeroOut tells the tunnel to accept zero-hop peers
func SetAllowZeroOut(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.OutAllowZeroHop = b
		return nil
	}
}

//SetCompress tells clients to use compression
func SetCompress(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.UseCompression = b
		return nil
	}
}

//SetFastRecieve tells clients to use compression
func SetFastRecieve(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.FastRecieve = b
		return nil
	}
}

//SetReduceIdle tells the connection to reduce it's tunnels during extended idle time.
func SetReduceIdle(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.ReduceIdle = b
		return nil
	}
}

//SetReduceIdleTime sets the time to wait before reducing tunnels to idle levels
func SetReduceIdleTime(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.ReduceIdleTime = 300000
		if u >= 6 {
			c.Conf.ReduceIdleTime = (u * 60) * 1000
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in minutes) %v", u)
	}
}

//SetReduceIdleTimeMs sets the time to wait before reducing tunnels to idle levels in milliseconds
func SetReduceIdleTimeMs(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.ReduceIdleTime = 300000
		if u >= 300000 {
			c.Conf.ReduceIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid reduce idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetReduceIdleQuantity sets minimum number of tunnels to reduce to during idle time
func SetReduceIdleQuantity(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if u < 5 {
			c.Conf.ReduceIdleQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid reduce tunnel quantity")
	}
}

//SetCloseIdle tells the connection to close it's tunnels during extended idle time.
func SetCloseIdle(b bool) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.CloseIdle = b
		return nil
	}
}

//SetCloseIdleTime sets the time to wait before closing tunnels to idle levels
func SetCloseIdleTime(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.CloseIdleTime = 300000
		if u >= 6 {
			c.Conf.CloseIdleTime = (u * 60) * 1000
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in minutes) %v", u)
	}
}

//SetCloseIdleTimeMs sets the time to wait before closing tunnels to idle levels in milliseconds
func SetCloseIdleTimeMs(u int) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		c.Conf.CloseIdleTime = 300000
		if u >= 300000 {
			c.Conf.CloseIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid close idle timeout(Measured in milliseconds) %v", u)
	}
}

//SetAccessListType tells the system to treat the accessList as a whitelist
func SetAccessListType(s string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if s == "whitelist" {
			c.Conf.AccessListType = "whitelist"
			return nil
		} else if s == "blacklist" {
			c.Conf.AccessListType = "blacklist"
			return nil
		} else if s == "none" {
			c.Conf.AccessListType = ""
			return nil
		} else if s == "" {
			c.Conf.AccessListType = ""
			return nil
		}
		return fmt.Errorf("Invalid Access list type(whitelist, blacklist, none)")
	}
}

//SetAccessList tells the system to treat the accessList as a whitelist
func SetAccessList(s []string) func(*SAMClientVPN) error {
	return func(c *SAMClientVPN) error {
		if len(s) > 0 {
			for _, a := range s {
				c.Conf.AccessList = append(c.Conf.AccessList, a)
			}
			return nil
		}
		return nil
	}
}
