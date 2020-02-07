module github.com/RTradeLtd/libanonvpn

go 1.12

require (
	crawshaw.io/littleboss v0.0.0-20190317185602-8957d0aedcce
	github.com/eyedeekay/accessregister v0.0.0-20190908214045-2f83369c289b
	github.com/eyedeekay/canal v0.0.26
	github.com/eyedeekay/checki2cp v0.0.0-20191027173419-138f1b4882b2
	github.com/eyedeekay/go-i2cp v0.0.0-20190711020517-c0bce4e7b750 // indirect
	github.com/eyedeekay/httptunnel v0.0.0-20191017011116-3b144b52941f
	github.com/eyedeekay/outproxy v0.0.0-20190908174238-22bd71d43733
	github.com/eyedeekay/portcheck v0.0.0-20190218044454-bb8718669680
	github.com/eyedeekay/sam-forwarder v0.32.1
	github.com/eyedeekay/sam3 v0.32.2
	github.com/eyedeekay/udptunnel v0.0.92
	github.com/justinas/nosurf v0.0.0-20190416172904-05988550ea18
	github.com/kardianos/service v1.0.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/zieckey/goini v0.0.0-20180118150432-0da17d361d26
	github.com/zserge/lorca v0.1.8
	github.com/zserge/webview v0.0.0-20190123072648-16c93bcaeaeb
	golang.org/x/crypto v0.0.0-20191122220453-ac88ee75c92c // indirect
	golang.org/x/mobile v0.0.0-20191002175909-6d0d39b2ca82
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933 // indirect
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // indirect
)

replace github.com/eyedeekay/gosam v0.1.1-0.20190705071001-d8c0f81c783e => github.com/eyedeekay/goSam v0.1.1-0.20190705071001-d8c0f81c783e

replace golang.org/x/lint v0.0.0 => github.com/golang/lint v0.0.0
