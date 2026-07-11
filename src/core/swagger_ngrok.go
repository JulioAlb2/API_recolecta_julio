package core

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	swaggerHostRe    = regexp.MustCompile(`"host"\s*:\s*"[^"]*"`)
	swaggerSchemesRe = regexp.MustCompile(`"schemes"\s*:\s*\[[^\]]*\]`)
)

// InjectNgrokSkipHeaderInSwagger:
// 1) Reescribe doc.json para que host/schemes coincidan con la URL real
//    (ngrok/IP) y Swagger deje de generar curl a http://localhost:8080.
// 2) Inyecta ngrok-skip-browser-warning en el HTML de Swagger UI.
//
// No abre CORS: solo corrige la documentación servida al cliente.
func InjectNgrokSkipHeaderInSwagger() gin.HandlerFunc {
	const snippet = `<script>
(function () {
  var originalFetch = window.fetch;
  window.fetch = function (input, init) {
    init = init || {};
    var headers = new Headers(init.headers || {});
    headers.set("ngrok-skip-browser-warning", "true");
    init.headers = headers;
    return originalFetch(input, init);
  };
})();
</script>`

	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api/swagger") {
			c.Next()
			return
		}

		writer := &swaggerBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		c.Next()

		contentType := writer.Header().Get("Content-Type")
		body := writer.body.Bytes()
		path := c.Request.URL.Path

		// doc.json / swagger JSON: forzar host de la petición entrante
		if strings.Contains(contentType, "json") || strings.HasSuffix(path, "doc.json") || strings.HasSuffix(path, "swagger.json") {
			host := c.Request.Host
			if host == "" {
				host = c.GetHeader("X-Forwarded-Host")
			}
			proto := c.GetHeader("X-Forwarded-Proto")
			if proto == "" {
				if c.Request.TLS != nil {
					proto = "https"
				} else if strings.Contains(host, "ngrok") {
					proto = "https"
				} else {
					proto = "http"
				}
			}

			if host != "" {
				jsonBody := string(body)
				jsonBody = swaggerHostRe.ReplaceAllString(jsonBody, `"host": "`+host+`"`)
				jsonBody = swaggerSchemesRe.ReplaceAllString(jsonBody, `"schemes": ["`+proto+`"]`)
				body = []byte(jsonBody)
				writer.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(body)))
				log.Printf("Swagger doc reescrito: host=%s scheme=%s", host, proto)
			}
		}

		if strings.Contains(contentType, "text/html") {
			html := string(body)
			if strings.Contains(html, "</body>") && !strings.Contains(html, "ngrok-skip-browser-warning") {
				html = strings.Replace(html, "</body>", snippet+"</body>", 1)
			}
			body = []byte(html)
			writer.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(body)))
		}

		_, _ = writer.ResponseWriter.Write(body)
	}
}

type swaggerBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *swaggerBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *swaggerBodyWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}
