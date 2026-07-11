package main

import (
	"strings"

	"github.com/vicpoo/API_recolecta/docs"
)

func configureSwaggerDocs() {
	tmpl := docs.SwaggerInfo.SwaggerTemplate
	if strings.Contains(tmpl, `"security":[{"BearerAuth":[]}]`) {
		return
	}

	tmpl = strings.Replace(tmpl, `"paths": {`, `"security":[{"BearerAuth":[]}],"paths":{`, 1)

	for _, ex := range []struct{ path, method string }{
		{"/api/empleados/login", "post"},
		{"/api/ciudadanos/login", "post"},
		{"/api/ciudadanos", "post"},
		{"/api/colonia", "get"},
		{"/api/colonia/{id}", "get"},
	} {
		tmpl = exemptPublicSwaggerEndpoint(tmpl, ex.path, ex.method)
	}

	docs.SwaggerInfo.SwaggerTemplate = tmpl
}

func exemptPublicSwaggerEndpoint(tmpl, path, method string) string {
	needle := `"` + path + `": {
            "` + method + `": {`
	replacement := `"` + path + `": {
            "` + method + `": {
                "security": [],`
	return strings.Replace(tmpl, needle, replacement, 1)
}
