package handlers

import (
	"net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
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
func Info(logger zerolog.Logger, resp dtos.StatusResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Msg("/api/v1/info")
		_ = server.JSON(w, http.StatusOK, resp)
	}
}
