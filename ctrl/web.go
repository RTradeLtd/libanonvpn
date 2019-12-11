package samtunnelhandler

import (
	"net/http"
)

import (
	samtunnelhandler "github.com/eyedeekay/sam-forwarder/handler"
	samtunnel "github.com/eyedeekay/sam-forwarder/interface"
)

type VPNTunnelHandler struct {
	*samtunnelhandler.TunnelHandler
}

func (t *VPNTunnelHandler) Printdivf(id, key, value string, rw http.ResponseWriter, req *http.Request) {
	t.TunnelHandler.Printdivf(id, key, value, rw, req)
}

func (t *VPNTunnelHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	t.TunnelHandler.ServeHTTP(rw, req)
}

func NewTunnelHandler(ob samtunnel.SAMTunnel, err error) (*VPNTunnelHandler, error) {
	var t VPNTunnelHandler
	t.TunnelHandler, err = samtunnelhandler.NewTunnelHandler(ob, err)
	return &t, err
}
