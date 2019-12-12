
go-anonvpn ([home](/))
======================

Library for providing and connecting to VPN's over the I2P network.
Daemon, web client, and terminal client. This is an automatically
configuring, automatically deploying, automatically multihopping
pseudonymous VPN.

Installation
-------------

For now, the recommended way to install is with ```go get```

       go get -u -d -tags cli github.com/RTradeLtd/libanonvpn/cmd/anonvpn

Development
-----------

  - To install the dependencies needed to build this locally run
   `make dependencies` which will install:
    - `libgtk-3-dev`
    - `libwebkit2gtk-4.0-dev`
    - `build-essential`
    - `gomobile` which can be installed with `go get golang.org/x/mobile/cmd/gomobile`
  - You will also need a valid installation of 'go 1.12.x'

Example Usage
-------------

### Server-Side

Start by creating a server configuration file like the one found in
/etc/anonvpn/anonvpn.ini. Then run the server using that file:

        ./anonvpn -file server.ini


### Client-Side

When the server is started, it will create a minimum viable configuration
file for clients. You can run with a similar command:

        ./anonvpn -file client.ini


```
Usage of ./cmd/anonvpn/anonvpn:
  -accesslist string
    	Type of access list to use, can be "whitelist" "blacklist" or "none". (default "none")
  -addr string
    	(client) IP address of virtual network interface (default "10.79.0.2")
  -bch
    	Use a bitcoin cash wallet.
  -btc
    	Use a bitcoin wallet, true by default(rationale=widely adopted). (default true)
  -chromeuser string
    	user to run Chrome as, usually your desktop user (default "idk")
  -client
    	Client mode(true or false). (default true)
  -clientconf string
    	(Server Only) Output a client config file to the specified path (default "client.ini")
  -closeidle
    	Close tunnel after idle for a specified time(true or false).
  -closeidletime int
    	Close tunnel group after X (milliseconds). (default 600000)
  -compression
    	Uze gzip(true or false).
  -css string
    	custom CSS for web interface (default "css/styles.css")
  -destination string
    	Destination to connect client's to by default.
  -directory string
    	Directory to save tunnel configuration file in.
  -encryptleaseset
    	Use an encrypted leaseset(true or false). (default true)
  -eth
    	Use an ethereum wallet.
  -file string
    	Use an ini file for configuration(config file options override passed arguments for now). (default "none")
  -hashhash string
    	32-word mnemonic representing a .b32.i2p address(will output .b32.i2p address and quit)
  -host string
    	(server) IP address of virtual network interface (default "10.79.0.1")
  -inbackups int
    	Set inbound tunnel backup quantity(0 to 5). (default 3)
  -inlength int
    	Set inbound tunnel length(0 to 7). (default 1)
  -inquantity int
    	Set inbound tunnel quantity(0 to 15). (default 5)
  -invariance int
    	Set inbound tunnel length variance(-7 to 7).
  -javascript string
    	custom JS for web interface (default "js/scripts.js")
  -k string
    	key for encrypted leaseset (default "none")
  -littleboss string
    	instruct the littleboss:
    	
    	start:		start and manage this process using service name "name"
    	stop:		signal the littleboss to shutdown the process
    	status:		print statistics about the running littleboss
    	reload:		restart the managed process using the executed binary
    	bypass:		disable littleboss, run the program directly (default "bypass")
  -ltc
    	Use a litecoin wallet.
  -mnemonic string
    	Load or restore a wallet from the mnemonic string(Must be quoted).
  -name string
    	Tunnel name, this must be unique but can be anything. (default "anonvpn")
  -outbackups int
    	Set outbound tunnel backup quantity(0 to 5). (default 3)
  -outlength int
    	Set outbound tunnel length(0 to 7). (default 1)
  -outproxy string
    	Configure a SOCKS outproxy with your wallet proxy(i2p mode)
  -outquantity int
    	Set outbound tunnel quantity(0 to 15). (default 5)
  -outvariance int
    	Set outbound tunnel length variance(-7 to 7).
  -password string
    	password for web admin panel
  -persistident
    	Use saved file and persist tunnel(If false, tunnel will not persist after program is stopped.
  -pk string
    	private key for encrypted leaseset (default "none")
  -proxy string
    	Proxy to use for the wallet connection.(Tor, i2p, or host:port). (default "i2p")
  -psk string
    	private signing key for encrypted leaseset (default "none")
  -rate int
    	Set a payment requirement to authorize new clients
  -reduceidle
    	Reduce tunnel quantity when idle for a specified time(true or false).
  -reduceidlequantity int
    	Reduce idle tunnel quantity to X (0 to 5). (default 3)
  -reduceidletime int
    	Reduce tunnel quantity after X (milliseconds). (default 600000)
  -requirepass string
    	Require a password to request service information.
  -samhost string
    	SAM host (default "127.0.0.1")
  -samport string
    	SAM port (default "7656")
  -signaturetype string
    	Signature type
  -skipi2cp
    	Skip I2CP Port check for standalone router
  -start
    	Start a tunnel with the passed parameters(Otherwise, they will be treated as default values). (default true)
  -username string
    	username for web admin panel (default "go-anonvpn")
  -wallet string
    	File to store the wallet in. (default "vpnwallet.dat")
  -walletnet string
    	Which network to use the wallet on (mainnet, testnet, regtest). (default "mainnet")
  -walletpass string
    	password to use for the wallet. (default "ChangeMe")
  -webface
    	Start web administration interface (default true)
  -webport string
    	Web interface port (default "7959")
  -zec
    	Use a zerocash wallet, true by default(rationale=mainstream privacy coin). (default true)
  -zeroin
    	Allow zero-hop, non-anonymous tunnels in(true or false).
  -zeroout
    	Allow zero-hop, non-anonymous tunnels out(true or false).
```


