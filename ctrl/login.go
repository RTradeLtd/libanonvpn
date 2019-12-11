package samtunnelhandler

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func (m *TunnelHandlerMux) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if creds.Username != m.user {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if creds.Password != m.password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	m.sessionToken, err = GenerateRandomString(32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   m.sessionToken,
		Expires: time.Now().Add(10 * time.Minute),
	})
}

func (m *TunnelHandlerMux) Home(w http.ResponseWriter, r *http.Request) {
	if m.CheckCookie(w, r) == false {
		return
	}
	//	fmt.Fprintf(w, "URL PATH", r.URL.Path)
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index.html", 301)
		fmt.Fprintf(w, "redirecting to index.html")
		return
	}
	r2, err := http.NewRequest("GET", r.URL.Path+"/color", r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<!DOCTYPE html>\n")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, "<head>\n")
	fmt.Fprintf(w, "  <link rel=\"stylesheet\" href=\"/styles.css\">")
	fmt.Fprintf(w, "</head>\n")
	fmt.Fprintf(w, "<body>\n")
	fmt.Fprintf(w, "<h1>\n")
	w.Write([]byte(fmt.Sprintf("<a href=\"/index.html\">Welcome %s! you are serving %d tunnels. </a>\n", m.user, len(m.tunnels))))
	fmt.Fprintf(w, "</h1>\n")
	fmt.Fprintf(w, "  <div id=\"toggleall\" class=\"global control\">\n")
	fmt.Fprintf(w, "    <a href=\"#\" onclick=\"toggle_visibility_class('%s');\">Show/Hide %s</a>\n", "prop", "all")
	fmt.Fprintf(w, "  </div>\n")
	for _, tunnel := range m.Tunnels() {
		tunnel.ServeHTTP(w, r2)
	}
	fmt.Fprintf(w, "  <script src=\"/scripts.js\"></script>\n")
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, "</html>\n")
}

func (m *TunnelHandlerMux) CSS(w http.ResponseWriter, r *http.Request) {
	if m.CheckCookie(w, r) == false {
		return
	}
	w.Header().Add("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s\n", m.cssString)))
}

func (m *TunnelHandlerMux) JS(w http.ResponseWriter, r *http.Request) {
	if m.CheckCookie(w, r) == false {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s\n", m.jsString)))
}
