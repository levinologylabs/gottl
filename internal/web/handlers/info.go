package handlers

import (
	"net/http"

	"github.com/hay-kot/httpkit/server"
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
//	@Router			/api/v1/status [GET]
func Info(resp dtos.StatusResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = server.JSON(w, http.StatusOK, resp)
	}
}
