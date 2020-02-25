export GOPRIVATE=github.com/RTradeLtd/*,github.com/eyedeekay/*
export GO111MODULE=on

all: fmt build nostatic android windowsall osx nsis checkinstall doc

modtidy:
	go mod tidy

tap:
	wget -O nsis/wintap.exe https://build.openvpn.net/downloads/releases/latest/tap-windows-latest-stable.exe

tun:
	wget -O nsis/wintun.msi https://github.com/meshsocket/WintunInstaller/releases/download/0.6/wintun-amd64-0.6.msi
	wget -O nsis/wintun32.msi https://github.com/meshsocket/WintunInstaller/releases/download/0.6/wintun-x86-0.6.msi

tuntap: tun tap

cli:
	GOARCH=amd64 go build -a -tags="$(gotags) netgo cli" \
		-ldflags '-w -extldflags "-static"' \
		-o cmd/anonvpn/anonvpn \
		./cmd/anonvpn

build: cli
	GOARCH=amd64 go build -a -tags="$(gotags) netgo static" \
		-ldflags '-w -extldflags "-static"' \
		-o cmd/anonvpn/anonvpn-gui \
		./cmd/anonvpn

nostatic:
	GOARCH=amd64 go build -a -tags="$(gotags) nostatic" \
		-o cmd/anonvpn/anonvpn-nostatic \
		./cmd/anonvpn

windowsall: windows windows32

windows: syso
	GOARCH=amd64 GOOS=windows go build -a -tags="$(gotags) netgo static" \
		-buildmode=exe \
		-o cmd/anonvpn/anonvpn.exe \
		./cmd/anonvpn
	rm -f cmd/anonvpn/out.syso

windows32: syso
	GOARCH=386 GOOS=windows go build -a -tags="$(gotags) netgo static" \
		-buildmode=exe \
		-o cmd/anonvpn/anonvpn-32.exe \
		./cmd/anonvpn
	rm -f cmd/anonvpn/out.syso

osx:
	GOARCH=amd64 GOOS=darwin go build -a -tags="$(gotags) netgo cli" \
		-o cmd/anonvpn/anonvpn-osx \
		./cmd/anonvpn

wasm:
	GOARCH=wasm GOOS=js go build -a -tags="$(gotags) netgo cli" \
		-o cmd/anonvpn/anonvpn.js \
		./cmd/anonvpn

js:
	gopherjs build --tags="$(gotags) netgo cli" \
		-o cmd/anonvpn/anonvpn.js \
		./cmd/anonvpn

java: js
	rhino cmd/anonvpn/anonvpn.js

android:
	GOOS=android gomobile build -target=android -tags="$(gotags) netgo cli android" \
		-o cmd/anonvpn/anonvpn.apk \
		./cmd/anonvpn/

delreseed:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r libanonvpn -t reseed; true

reseed: installer updater
	rm -f etc/anonvpn/reseed.zip
	wget -O etc/anonvpn/reseed.zip http://localhost:7657/createreseed
	gothub release -s $(GITHUB_TOKEN) -p -u $(USER_GH) -r libanonvpn -t reseed -d "Privacy-Enhanced VPN - $(LATEST_DESC)"
	gothub upload -s $(GITHUB_TOKEN) -f "etc/anonvpn/reseed.zip" -n "reseed.zip" -u $(USER_GH) -r libanonvpn -t reseed -l "Reseed File" -R
	gothub upload -s $(GITHUB_TOKEN) -f "i2pinstall.exe" -n "i2pinstall.exe" -u $(USER_GH) -r libanonvpn -t reseed -l "I2P Dev Build" -R
	gothub upload -s $(GITHUB_TOKEN) -f "i2pupdate.zip" -n "i2pupdate.zip" -u $(USER_GH) -r libanonvpn -t reseed -l "I2P Dev Build updater.zip" -R

syso:
	syso -o cmd/anonvpn/out.syso

zip:

winlicense:
	cat client/LICENSE.md server/LICENSE.md | unix2dos | tee LICENSES.txt

nsis:geti2p windows windows-32 winstall
winstall: winlicense tuntap
	makensis nsis/installer.nsi

-include nsis/geti2p.mk

fmt: go-fmt

deps:
	go get -u github.com/RTradeLtd/libanonvpn/cmd/anonvpn

setcap:
	sudo setcap "cap_net_admin+eip cap_net_bind_service+eip cap_net_raw+eip" /usr/bin/anonvpn
	sudo getcap /usr/bin/anonvpn

install:
	install -m755 ./cmd/anonvpn/anonvpn /usr/bin/anonvpn
	install -m755 ./cmd/anonvpn/anonvpn-gui /usr/bin/anonvpn-gui
	mkdir -p /etc/anonvpn
	install -m644 ./etc/anonvpn/anonvpn-server.ini /etc/anonvpn/anonvpn-server.ini
	install -m644 ./etc/anonvpn/anonvpn.ini /etc/anonvpn/anonvpn.ini
	install -m644 ./etc/anonvpn/i2cp.conf /etc/anonvpn/anonvpn.conf

try-server:
	sudo -b ./cmd/anonvpn/anonvpn -file etc/anonvpn/anonvpn-server.ini 2>&1 | tee server.log

try:
	sudo -b ./cmd/anonvpn/anonvpn -file etc/anonvpn/anonvpn-server.ini 2>&1 | tee server.log

try-client:
	sudo -b ./cmd/anonvpn/anonvpn -file cvpnserver.ini 2>&1 | tee server.log

try-rtrade:
	sudo -b ./cmd/anonvpn/anonvpn -file rtrade-testserver.ini 2>&1 | tee server.log


try-server-osx:
	sudo -b ./cmd/anonvpn/anonvpn-osx -file etc/anonvpn/anonvpn-server.ini 2>&1 | tee server.log

try-client-osx:
	sudo -b ./cmd/anonvpn/anonvpn-osx -file cvpnserver.ini 2>&1 | tee server.log

clean: fmt clean-pkg
	find . -name '*.i2pkeys' -exec rm -v {} \;
	find . -name '*.log' -exec rm -v {} \;
	find . -name '*.syso' -exec rm -v {} \;
	rm -vfr *.exe .geti2p.url backup*.tgz description-pak libanonvpn
	rm -f cmd/anonvpn/anonvpn*
	go clean ./cmd/anonvpn

doc: head example help
	grep -v ':::' README.0.md > README.md

head:
	@echo "" | tee README.0.md
	@echo "::: {.content .toplevel}" | tee -a README.0.md
	@echo "libanonvpn ([home](/))" | tee -a README.0.md
	@echo "======================" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "Library for providing and connecting to VPN's over the I2P network." | tee -a README.0.md
	@echo "Daemon, web client, and terminal client. This is an automatically" | tee -a README.0.md
	@echo "configuring, automatically deploying, automatically multihopping" | tee -a README.0.md
	@echo "pseudonymous VPN." | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "::: {.content .installsource}" | tee -a README.0.md
	@echo "Installation" | tee -a README.0.md
	@echo "-------------" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "For now, the recommended way to install is with \`\`\`go get\`\`\`" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo '       go get -u -d -tags cli github.com/RTradeLtd/libanonvpn/cmd/anonvpn' | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md

about:
	@echo "::: {.content .privacy}" | tee -a README.0.md
	@echo "Turn-Key Privacy Advantages" | tee -a README.0.md
	@echo "---------------------------" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "  - **Trustless Multihopping:** libanonvpn uses the I2P network to negotiate" | tee -a README.0.md
	@echo "   encrypted, peer-to-peer tunnels which are used to establish the VPN hops" | tee -a README.0.md
	@echo "   without revealing the true IP address of any tunnel participant to any" | tee -a README.0.md
	@echo "   other participant.This is in contrast to Wireguard, where you can multihop," | tee -a README.0.md
	@echo "   but the configuration requires you to coordinate with multiple providers," | tee -a README.0.md
	@echo "   we can simply do it automatically." | tee -a README.0.md
	@echo "  - **Pseudonymous Subscription:** libanonvpn integrates a multi-wallet which" | tee -a README.0.md
	@echo "   is your virtual subscription. Pay for the VPN as-you-go by transferring" | tee -a README.0.md
	@echo "   money to the VPN wallet, and use a method similar to electrum's \"Change\"" | tee -a README.0.md
	@echo "   addresses to separate your payment identities from the rest of your crypto" | tee -a README.0.md
	@echo "   usage. This is so we can't correlate your specifidc traffic to your payment." | tee -a README.0.md
	@echo "   Also, pay as you go. Your subscription is the content of your wallet, never" | tee -a README.0.md
	@echo "   pay for more time than you need." | tee -a README.0.md
	@echo "  - **Resistant to Analysis:** libanonvpn inherits the traffic-obfuscation and" | tee -a README.0.md
	@echo "   encryption properties of the underlying I2P network, making it nearly" | tee -a README.0.md
	@echo "   impossible to identify and block usage of the VPN. With the assistance" | tee -a README.0.md
	@echo "   of \"helper\" applications, it can be impossible to block on even the" | tee -a README.0.md
	@echo "   most restricted networks." | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "::: {.content .others}" | tee -a README.0.md
	@echo "Other advantages" | tee -a README.0.md
	@echo "----------------" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "   - **Peer-to-Peer:** Your VPN server has one reachable name, and will only" | tee -a README.0.md
	@echo "    ever have one reachable name, which is cryptographically secure and " | tee -a README.0.md
	@echo "    addressable within the I2P network only. No one can ever impersonate your" | tee -a README.0.md
	@echo "    server unless they steal your private keys." | tee -a README.0.md
	@echo "   - **Self-Hosting:** A VPN Service can be hosted by anyone, anywhere, while" | tee -a README.0.md
	@echo "    maintaining a trust-less system, either gratis as one may in the interest of" | tee -a README.0.md
	@echo "    the public good or as a participant in a de-centralized VPN marketplace" | tee -a README.0.md
	@echo "    where servers can set their own prices based on the features and bandwidth" | tee -a README.0.md
	@echo "    or anonymity features they can provide." | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "::: {.content .planned}" | tee -a README.0.md
	@echo "### Planned Advantages" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "*These features are incomplete and non-operational. They will be moved as they*" | tee -a README.0.md
	@echo "*are implemented.*" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "   - **De-Centralization:** Soon, instead of subscribing to a fixed VPN endpoint" | tee -a README.0.md
	@echo "    one may query a de-centralized pool of nodes providing the services they" | tee -a README.0.md
	@echo "    require, in order to enhance their anonymity by making their exit to the" | tee -a README.0.md
	@echo "    internet difficult to predict. This pool of nodes will be anonymous and" | tee -a README.0.md
	@echo "    uncensorable.Some possible candidates for this have already been" | tee -a README.0.md
	@echo "    implemented, none decided on, most not mutually exclusive." | tee -a README.0.md
	@echo "   - **Server-Blinding:** Server-to-Server agreements, financial or otherwise," | tee -a README.0.md
	@echo "    may be used to further pool and obfuscate the origin of a connection, in order" | tee -a README.0.md
	@echo "    to protect server operators from persecution intended to target their" | tee -a README.0.md
	@echo "    subscribers. From the outside, it will be nearly impossible to determine" | tee -a README.0.md
	@echo "    which customer ID's are being served by which server ID's in real time" | tee -a README.0.md
	@echo "    or retro-actively." | tee -a README.0.md
	@echo "   - **Redistributable** Since this will work best with a diverse userbase, we" | tee -a README.0.md
	@echo "    want to make it as easy as possible to share this application safely with" | tee -a README.0.md
	@echo "    a local friend. The application will be capable of doing things like setting" | tee -a README.0.md
	@echo "    up a local F-Droid or Aptly repository for sharing the package on a LAN, or" | tee -a README.0.md
	@echo "    to expose the repository as an I2P-only site." | tee -a README.0.md
	@echo "    - **Bittorrent-over-I2P based Updates:** In order to ensure that the VPN is" | tee -a README.0.md
	@echo "    always up-to-date and accessible to people who already have it, a" | tee -a README.0.md
	@echo "    Bittorrent-over-I2P based updating system will be implemented where all" | tee -a README.0.md
	@echo "    participants will potentially provide updates to eachother from multiple" | tee -a README.0.md
	@echo "    anonymous sources." | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md

help:
	@echo "::: {.content .usage .help}" | tee -a README.0.md
	@echo '```' | tee -a README.0.md
	./cmd/anonvpn/anonvpn -h 2>&1 | tee -a README.0.md
	@echo '```' | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md
	@echo "" | tee -a README.0.md

example:
	@echo "::: {.content .usage .example}" | tee -a README.0.md
	@echo "Example Usage" | tee -a README.0.md
	@echo "-------------" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::: {.content .usage .server}" | tee -a README.0.md
	@echo "### Server-Side" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "Start by creating a server configuration file like the one found in" | tee -a README.0.md
	@echo "/etc/anonvpn/anonvpn.ini. Then run the server using that file:" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "        ./anonvpn -file server.ini" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "::::" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::: {.content .usage.client}" | tee -a README.0.md
	@echo "### Client-Side" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "When the server is started, it will create a minimum viable configuration" | tee -a README.0.md
	@echo "file for clients. You can run with a similar command:" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "        ./anonvpn -file client.ini" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo "::::" | tee -a README.0.md
	@echo "" | tee -a README.0.md
	@echo ":::" | tee -a README.0.md

pdf:
	pandoc --highlight-style=tango -f gfm README.md -t html5 -o README.pdf

all: fmt build doc

docker-server: build
	docker build -f server.Dockerfile -t eyedeekay/libanonvpn:server .

clean-server:
	docker rm -f libanonvpn-server; true

run-server: clean-server
	docker run -it \
		--net=anonvpni2p \
		--cap-add NET_ADMIN \
		--cap-add NET_RAW \
		--cap-add NET_BIND_SERVICE \
		--name libanonvpn-server \
		--privileged \
		-p 127.0.0.1:7959:7959 \
		eyedeekay/libanonvpn:server

docker: docker-server run-server

docker-client: build
	docker build -f client.Dockerfile -t eyedeekay/libanonvpn:client .

run-client:
	docker run --rm -it \
		--net=host \
		--cap-add NET_ADMIN \
		--cap-add NET_RAW \
		--cap-add NET_BIND_SERVICE \
		--name libanonvpn-client \
		--privileged \
		-p 127.0.0.1:7959:7959 \
		 eyedeekay/libanonvpn:client

# run standard go tooling for better readability
.PHONY: go-tidy
go-tidy: go-imports go-fmt
	go vet ./...
	golint ./...

# automatically add missing imports
.PHONY: go-imports
go-imports:
	find . -type f -path ./vendor -prune -o -name '*.go' -exec goimports -w {} \;

# format code and simplify if possible
.PHONY: go-fmt
go-fmt:
	find . -type f -path ./vendor -prune -o -name '*.go' -exec gofmt -s -w {} \;

.PHONY: dependencies
dependencies:
	sudo apt install dos2unix -y
	sudo apt install nsis -y
	sudo apt install rasqal -y
	sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev build-essential -y
	go get -u github.com/hallazzang/syso/...
	go get golang.org/x/mobile/cmd/gomobile

checkinstall:
	checkinstall --type=debian \
		--install=no \
		--fstrans=yes \
		--pkgname=libanonvpn \
		--pkgversion=$(VERSION) \
		--pkgrelease=testing \
		--pkggroup=net \
		--pakdir=./ \
		--maintainer=hankhill19580@gmail.com \
		--requires="i2p | i2pd, coreutils" \
		--nodoc \
		--deldoc=yes \
		--deldesc=yes \
		--delspec=yes \
		--default

docker-clientfile:
	docker cp libanonvpn-server:/opt/go/src/github.com/RTradeLtd/libanonvpn/cvpnserver.ini .

clean-tunsocks:
	rm -rf tunsocks

socks: clean-tunsocks
	git clone https://github.com/russdill/tunsocks
	cd tunsocks && \
		git submodule init && \
		git submodule update && \
		./autogen.sh && \
		./configure && \
		make

DEMOVERSION=0.32.01
VERSION=0.32.01
#USER_GH="eyedeekay"
USER_GH="RTradeLtd"
export ANDROID_NDK_HOME=$(HOME)/Workspace/android-ndk-r19c

clean-pkg:
	rm -f *.tgz *.tar.gz *.exe *.deb

tarball: all
	rm -f ../libanonvpn_*.tar.gz
	tar --exclude=".git" --exclude=backup*.tgz -czf ../libanonvpn_$(VERSION).tar.gz .
	mv ../libanonvpn_$(VERSION).tar.gz .

STABLE_DESC=This release has undergone testing by the developers and is recommended for most users. It is always a copy of the most recent tagged release.
LATEST_DESC=This release is always built from the latest buildable code and may contain bugs.

version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r go-anonvpn -t v$(VERSION) -d "Privacy-Enhanced VPN"

latest:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r go-anonvpn -t latest; true
	gothub release -s $(GITHUB_TOKEN) -p -u $(USER_GH) -r go-anonvpn -t latest -d "Privacy-Enhanced VPN - $(LATEST_DESC)"

tag: version latest

pages: docs
	cd docs && make pages

host:
	cd docs && make host

sed:
	find . -name '*.go' -exec sed -i 's|github.com/RTradeLtd/libanonvpn/wallet|github.com/RTradeLtd/davpn-gateways/wallet|g' {} \;

find:
	find . -name '*.go' -exec grep canal {} \;