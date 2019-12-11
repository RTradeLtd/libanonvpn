package samtunnelhandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TunnelHandlerMux struct {
	http.Server
	pagenames    []string
	tunnels      []*VPNTunnelHandler
	user         string
	password     string
	sessionToken string
	cssString    string
	jsString     string
}

func (m *TunnelHandlerMux) ListenAndServe() {
	m.Server.ListenAndServe()
}

func (m *TunnelHandlerMux) CheckCookie(w http.ResponseWriter, r *http.Request) bool {
	if m.password != "" {
		if m.sessionToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return false
			}
			w.WriteHeader(http.StatusBadRequest)
			return false
		}
		if m.sessionToken != c.Value {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
	}
	return true

}

func (m *TunnelHandlerMux) HandlerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.CheckCookie(w, r) == false {
			return
		}
		m.ColorHeader(h, r, w)
	})
}

func (t *TunnelHandlerMux) Tunnels() []*VPNTunnelHandler {
	return t.tunnels
}

func (t *TunnelHandlerMux) ClientTunnels() []*VPNTunnelHandler {
	var tunnels []*VPNTunnelHandler
	for _, v := range t.tunnels {
		if v.GetType() == "vpnclient" {
			tunnels = append(tunnels, v)
		}
	}
	return tunnels
}

func (t *TunnelHandlerMux) ServerTunnels() []*VPNTunnelHandler {
	var tunnels []*VPNTunnelHandler
	for _, v := range t.tunnels {
		if v.GetType() == "vpnserver" {
			tunnels = append(tunnels, v)
		}
	}
	return tunnels
}

func (m *TunnelHandlerMux) Append(v *VPNTunnelHandler) *TunnelHandlerMux {
	if m == nil {
		return m
	}
	for _, prev := range m.tunnels {
		if v.ID() == prev.ID() {
			log.Printf("v.ID() found, %s == %s", v.ID(), prev.ID())
			return m
		}
	}
	log.Printf("Adding tunnel ID: %s", v.ID())
	m.tunnels = append(m.tunnels, v)
	Handler := m.Handler.(*http.ServeMux)
	Handler.Handle(fmt.Sprintf("/%d", len(m.tunnels)), m.HandlerWrapper(v))
	Handler.Handle(fmt.Sprintf("/%s", v.ID()), m.HandlerWrapper(v))
	Handler.Handle(fmt.Sprintf("/%d/color", len(m.tunnels)), m.HandlerWrapper(v))
	Handler.Handle(fmt.Sprintf("/%s/color", v.ID()), m.HandlerWrapper(v))
	m.Handler = Handler
	return m
}

func ReadFile(filename string) (string, error) {
	r, e := ioutil.ReadFile(filename)
	return string(r), e
}

func NewTunnelHandlerMux(host, port, user, password, css, javascript string) *TunnelHandlerMux {
	var m TunnelHandlerMux
	m.Addr = host + ":" + port
	Handler := http.NewServeMux()
	m.pagenames = []string{"index.html", "index", ""}
	m.user = user
	m.password = password
	m.sessionToken = ""
	m.tunnels = []*VPNTunnelHandler{}
	var err error
	m.cssString, err = ReadFile(css)
	if err != nil {
		m.cssString = DefaultCSS()
	}
	m.jsString, err = ReadFile(javascript)
	if err != nil {
		m.jsString = DefaultJS()
	}
	for _, v := range m.pagenames {
		Handler.HandleFunc(fmt.Sprintf("/%s", v), m.Home)
	}
	Handler.HandleFunc("/styles.css", m.CSS)
	Handler.HandleFunc("/scripts.js", m.JS)
	if m.password != "" {
		Handler.HandleFunc("/login", m.Signin)
	}
	m.Handler = Handler
	return &m
}
