package docs

import (
	_ "embed"
	"net/http"
)

//go:embed apidocs.swagger.json
var specJSON []byte

// SpecPath is where the raw OpenAPI JSON is served.
const SpecPath = "/swagger/apidocs.swagger.json"

// swaggerUI loads Swagger UI from a CDN and points it at the embedded spec.
const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>rbac-go API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({ url: "./apidocs.swagger.json", dom_id: "#swagger-ui" });
    };
  </script>
</body>
</html>`

// Handler serves the OpenAPI JSON and the Swagger UI page under
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(SpecPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(specJSON)
	})
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerUI))
	})
	return mux
}
