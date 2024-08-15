package handlers

import (
	"net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web/extractors"
)

type AdminController struct {
	adminservice *services.AdminService
}

func NewAdminController(adminservice *services.AdminService) *AdminController {
	return &AdminController{adminservice: adminservice}
}

// GetAllUsers godoc
//
//	@Tags			Admin
//	@Summary		Get all users
//	@Description	Get all users
//	@Accept			json
//	@Produce		json
//	@Param			page	query		dtos.Pagination	true	"Pagination details"
//	@Success		200		{object}	dtos.PaginationResponse[dtos.User]
//	@Router			/api/v1/admin/users [GET]
func (ac *AdminController) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	page, err := extractors.Query[dtos.Pagination](r)
	if err != nil {
		return err
	}

	users, err := ac.adminservice.GetAllUsers(r.Context(), page.WithDefaults())
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, users)
}
