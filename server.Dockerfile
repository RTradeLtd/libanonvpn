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
#RUN GO111MODULE=on go mod vendor
#RUN make build
RUN make install
#RUN apk del git make
CMD ./cmd/bin/anonvpn \
    -samhost=i2p \
    -samport=7657 \
    -littleboss=start \
    -file /opt/go/src/github.com/RTradeLtd/libanonvpn/etc/anonvpn/anonvpn-server.ini

