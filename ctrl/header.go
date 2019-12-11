package samtunnelhandler

import (
	"fmt"
	"net/http"
	"strings"
)

func (m *TunnelHandlerMux) ColorHeader(h http.Handler, r *http.Request, w http.ResponseWriter) {
	if !strings.HasSuffix(r.URL.Path, "color") {
		h.ServeHTTP(w, r)
	} else {
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
		h.ServeHTTP(w, r)
		fmt.Fprintf(w, "  <script src=\"/scripts.js\"></script>\n")
		fmt.Fprintf(w, "</body>\n")
		fmt.Fprintf(w, "</html>\n")
	}
}
