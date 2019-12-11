package samforwardervpnserver

/*
import (
	"fmt"
	"net/url"

	"golang.org/x/net/proxy"

	wallet "github.com/OpenBazaar/wallet-interface"
	"github.com/btcsuite/btcd/chaincfg"
)

func SetWalletFilePath(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.InnerWallet.WalletFilePath = s
		return nil
	}
}

func SetWalletPassword(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.InnerWallet.ChAngeMe = s
		return nil
	}
}

func SetChainConfigType(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if s == "mainnet" {
			c.InnerWallet.ChainCfg = &chaincfg.MainNetParams
		} else if s == "testnet" {
			c.InnerWallet.ChainCfg = &chaincfg.TestNet3Params
		} else if s == "regtest" {
			c.InnerWallet.ChainCfg = &chaincfg.RegressionNetParams
		} else {
			return fmt.Errorf("Invalid parameter: (Must be mainnet, testnet, regtest)")
		}
		return nil
	}
}

func SetUseBitcoin(b bool) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if c.InnerWallet.CoinTypes == nil {
			c.InnerWallet.CoinTypes = make(map[wallet.CoinType]bool)
		}
		c.InnerWallet.CoinTypes[wallet.Bitcoin] = b
		return nil
	}
}

func SetUseBitcoinCash(b bool) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if c.InnerWallet.CoinTypes == nil {
			c.InnerWallet.CoinTypes = make(map[wallet.CoinType]bool)
		}
		c.InnerWallet.CoinTypes[wallet.BitcoinCash] = b
		return nil
	}
}

func SetUseLitecoin(b bool) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if c.InnerWallet.CoinTypes == nil {
			c.InnerWallet.CoinTypes = make(map[wallet.CoinType]bool)
		}
		c.InnerWallet.CoinTypes[wallet.Litecoin] = b
		return nil
	}
}

func SetUseZerocash(b bool) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if c.InnerWallet.CoinTypes == nil {
			c.InnerWallet.CoinTypes = make(map[wallet.CoinType]bool)
		}
		c.InnerWallet.CoinTypes[wallet.Zcash] = b
		return nil
	}
}

func SetUseEthereum(b bool) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if c.InnerWallet.CoinTypes == nil {
			c.InnerWallet.CoinTypes = make(map[wallet.CoinType]bool)
		}
		c.InnerWallet.CoinTypes[wallet.Ethereum] = b
		return nil
	}
}

func SetProxyType(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if s == "tor" {
			u, err := url.Parse("socks://127.0.0.1:9050")
			c.InnerWallet.MultiWalletConf.Proxy, err = proxy.FromURL(u, proxy.Direct)
			return err
		} else if s == "i2p" {
			u, err := url.Parse("http://127.0.0.1:4444")
			c.InnerWallet.MultiWalletConf.Proxy, err = proxy.FromURL(u, proxy.Direct)
			return err
		} else if s != "" {
			u, err := url.Parse(s)
			c.InnerWallet.MultiWalletConf.Proxy, err = proxy.FromURL(u, proxy.Direct)
			return err
		}
		return nil
	}
}

func SetMnemonicLoad(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		if s != "" {
			c.InnerWallet.MultiWalletConf.Mnemonic = s
		}
		return nil
	}
}

func SetWalletBase32(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.AdvertiseBase32 = s
		return nil
	}
}

func SetWalletServicePrice(s int) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.AdvertiseServicePrice = s
		return nil
	}
}

func SetWalletRequirePassword(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.Password = s
		return nil
	}
}

func SetWalletTunName(s string) func(*SAMClientServerVPN) error {
	return func(c *SAMClientServerVPN) error {
		c.Conf.TunName = s
		return nil
	}
}
*/
