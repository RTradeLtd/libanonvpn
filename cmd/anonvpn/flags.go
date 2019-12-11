package main

import (
	"flag"
)

// Web Interface Options
var (
	webAdmin = flag.Bool("webface", true,
		"Start web administration interface")
	webPort = flag.String("webport", "7959",
		"Web interface port")
	webCSS = flag.String("css", "css/styles.css",
		"custom CSS for web interface")
	webJS = flag.String("javascript", "js/scripts.js",
		"custom JS for web interface")
	userName = flag.String("username", "go-anonvpn",
		"username for web admin panel")
	password = flag.String("password", "",
		"password for web admin panel")
	chromeUser = flag.String("chromeuser", runningUser,
		"user to run Chrome as, usually your desktop user")
)

// Application/VPN Options
var (
	iniFile = flag.String("file", "none",
		"Use an ini file for configuration(config file options override passed arguments for now).")
	targetHost = flag.String("host", "10.79.0.1",
		"(server) IP address of virtual network interface")
	samHost = flag.String("samhost", "127.0.0.1",
		"SAM host")
	samPort = flag.String("samport", "7656",
		"SAM port")
	tunnelHost = flag.String("addr", "10.79.0.2",
		"(client) IP address of virtual network interface")
	saveFile = flag.Bool("persistident", false,
		"Use saved file and persist tunnel(If false, tunnel will not persist after program is stopped.")
	clientConfig = flag.String("clientconf", "client.ini",
		"(Server Only) Output a client config file to the specified path")
	startUp = flag.Bool("start", true,
		"Start a tunnel with the passed parameters(Otherwise, they will be treated as default values).")
)

// I2P Options
var (
	encryptLeaseSet = flag.Bool("encryptleaseset", true,
		"Use an encrypted leaseset(true or false).")
	inAllowZeroHop = flag.Bool("zeroin", false,
		"Allow zero-hop, non-anonymous tunnels in(true or false).")
	outAllowZeroHop = flag.Bool("zeroout", false,
		"Allow zero-hop, non-anonymous tunnels out(true or false).")
	useCompression = flag.Bool("compression", false,
		"Uze gzip(true or false).")
	reduceIdle = flag.Bool("reduceidle", false,
		"Reduce tunnel quantity when idle for a specified time(true or false).")
	closeIdle = flag.Bool("closeidle", false,
		"Close tunnel after idle for a specified time(true or false).")
	client = flag.Bool("client", true,
		"Client mode(true or false).")
	/*server = flag.Bool("server", true,
	"Server mode(true or false).")*/
	peoplehash = flag.String("hashhash", "",
		"32-word mnemonic representing a .b32.i2p address(will output .b32.i2p address and quit)")
	sigType = flag.String("signaturetype", "",
		"Signature type")
	leaseSetKey = flag.String("k", "none",
		"key for encrypted leaseset")
	leaseSetPrivateKey = flag.String("pk", "none",
		"private key for encrypted leaseset")
	leaseSetPrivateSigningKey = flag.String("psk", "none",
		"private signing key for encrypted leaseset")
	targetDir = flag.String("directory", "",
		"Directory to save tunnel configuration file in.")
	targetDest = flag.String("destination", "",
		"Destination to connect client's to by default.")
	tunName = flag.String("name", "anonvpn",
		"Tunnel name, this must be unique but can be anything.")
	accessListType = flag.String("accesslist", "none",
		"Type of access list to use, can be \"whitelist\" \"blacklist\" or \"none\".")
	inLength = flag.Int("inlength", 1,
		"Set inbound tunnel length(0 to 7).")
	outLength = flag.Int("outlength", 1,
		"Set outbound tunnel length(0 to 7).")
	inQuantity = flag.Int("inquantity", 5,
		"Set inbound tunnel quantity(0 to 15).")
	outQuantity = flag.Int("outquantity", 5,
		"Set outbound tunnel quantity(0 to 15).")
	inVariance = flag.Int("invariance", 0,
		"Set inbound tunnel length variance(-7 to 7).")
	outVariance = flag.Int("outvariance", 0,
		"Set outbound tunnel length variance(-7 to 7).")
	inBackupQuantity = flag.Int("inbackups", 3,
		"Set inbound tunnel backup quantity(0 to 5).")
	outBackupQuantity = flag.Int("outbackups", 3,
		"Set outbound tunnel backup quantity(0 to 5).")
	reduceIdleTime = flag.Int("reduceidletime", 600000,
		"Reduce tunnel quantity after X (milliseconds).")
	closeIdleTime = flag.Int("closeidletime", 600000,
		"Close tunnel group after X (milliseconds).")
	reduceIdleQuantity = flag.Int("reduceidlequantity", 3,
		"Reduce idle tunnel quantity to X (0 to 5).")
	skipi2cp = flag.Bool("skipi2cp", false,
		"Skip I2CP Port check for standalone router")
)

// Wallet Options
var (
	walletFile = flag.String("wallet", "vpnwallet.dat",
		"File to store the wallet in.")
	walletPass = flag.String("walletpass", "ChangeMe",
		"password to use for the wallet.")
	mnemonicLoad = flag.String("mnemonic", "",
		"Load or restore a wallet from the mnemonic string(Must be quoted).")
	walletProxy = flag.String("proxy", "i2p",
		"Proxy to use for the wallet connection.(Tor, i2p, or host:port).")
	networkType = flag.String("walletnet", "mainnet",
		"Which network to use the wallet on (mainnet, testnet, regtest).")
	bitcoin = flag.Bool("btc", true,
		"Use a bitcoin wallet, true by default(rationale=widely adopted).")
	bitcoinCash = flag.Bool("bch", false,
		"Use a bitcoin cash wallet.")
	zeroCash = flag.Bool("zec", true,
		"Use a zerocash wallet, true by default(rationale=mainstream privacy coin).")
	litecoin = flag.Bool("ltc", false,
		"Use a litecoin wallet.")
	ethereum = flag.Bool("eth", false,
		"Use an ethereum wallet.")
	serviceprice = flag.Int("rate", 0,
		"Set a payment requirement to authorize new clients")
	requirepass = flag.String("requirepass", "",
		"Require a password to request service information.")
	outProxy = flag.String("outproxy", "",
		"Configure a SOCKS outproxy with your wallet proxy(i2p mode)")
)
