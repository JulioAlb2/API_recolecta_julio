package core

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// InjectNgrokSkipHeaderInSwagger reescribe el HTML de Swagger UI para que
// cada fetch envíe ngrok-skip-browser-warning. Sin esto, el plan gratis de
// ngrok intercepta el "Try it out" y Swagger muestra "Failed to fetch" / CORS.
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
