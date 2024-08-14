// Package docs servces the swagger.json file.
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed swagger.json
var swagger []byte

func SwaggerJSON(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(swagger)
	return err
}
