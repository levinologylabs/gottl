package docs

import (
	_ "embed"
	"net/http"
)

//go:embed swagger.json
var swagger []byte

func SwaggerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(swagger)
}
