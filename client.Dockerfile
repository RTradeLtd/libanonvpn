FROM alpine:3.10
ENV I2CP_HOME=/opt/go/src/github.com/RTradeLtd/libanonvpn/cmd/anonvpn
ENV GO_I2CP_CONF=/i2cp.docker.conf
ENV GOPATH=/opt/go/
RUN apk update
RUN apk add go git make musl-dev
RUN adduser -h /home/anonvpn -g 'anonvpn,,,,' -S -D anonvpn
COPY . /opt/go/src/github.com/RTradeLtd/libanonvpn
COPY etc/anonvpn/.i2cp.docker.conf \
    /opt/go/src/github.com/RTradeLtd/libanonvpn/etc/anonvpn/i2cp.docker.conf
WORKDIR /opt/go/src/github.com/RTradeLtd/libanonvpn
RUN GO111MODULE=on go mod vendor
RUN go version
RUN make build
RUN apk del git make
COPY cvpnserver.ini /opt/go/src/github.com/RTradeLtd/libanonvpn/etc/anonvpn/anonvpn-client.ini
CMD /opt/go/src/github.com/RTradeLtd/libanonvpn/cmd/anonvpn/anonvpn-webview \
    -littleboss=start \
    -file /opt/go/src/github.com/RTradeLtd/libanonvpn/etc/anonvpn/anonvpn-client.ini
