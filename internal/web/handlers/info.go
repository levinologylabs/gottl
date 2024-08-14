package handlers

import (
	"net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/data/dtos"
)

// Info godoc
//
//	@Summary		Get the status of the service
//	@Description	Get the status of the service
//	@Tags			status
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.StatusResponse
//	@Router			/api/v1/info [GET]
func Info(resp dtos.StatusResponse) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		return server.JSON(w, http.StatusOK, resp)
	}
}
