package handlers

import (
  "net/http"

	"github.com/jalevin/gottl/internal/core/server"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/internal/web/extractors"
)

type {{ .Computed.domain_var }}Controller struct {
  service *services.{{ .Computed.domain_var }}Service
}

func New{{ .Computed.domain_var }}Controller(service *services.{{ .Computed.domain_var }}Service) *{{ .Computed.domain_var }}Controller {
  return &{{ .Computed.domain_var }}Controller{
    service: service,
  }
}

// GetAll godoc
//
// @Tags      {{ .Computed.domain_var }}
// @Summary   List all {{ .Computed.domain_var }}s
// @Description List all {{ .Computed.domain_var }}s
// @Accept     json
// @Produce    json
// @Success    200     {object}  dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}]  "A list of {{ .Computed.domain_var }}s"
// @Router     /api/v1/{{ .Computed.domain_kebab }}s [GET]
func (uc *{{ .Computed.domain_var }}Controller) GetAll(w http.ResponseWriter, r *http.Request) error {
  page, err := extractors.Query[dtos.Pagination](r)
  if err != nil {
    return err
  }

  entities, err := uc.service.GetAll(r.Context(), page.WithDefaults())
  if err != nil {
    return err
  }

  return server.JSON(w, http.StatusOK, entities)
}

// Get godoc
//
// @Tags      {{ .Computed.domain_var }}
// @Summary   Get a {{ .Computed.domain_var }}
// @Description Get a {{ .Computed.domain_var }}
// @Accept     json
// @Produce    json
// @Param      id      path    string  true  "The {{ .Computed.domain_var }} ID"
// @Success    200     {object}  dtos.{{ .Computed.domain_var }}
// @Router     /api/v1/{{ .Computed.domain_kebab }}s/{id} [GET]
func (uc *{{ .Computed.domain_var }}Controller) Get(w http.ResponseWriter, r *http.Request) error {
  id, err := extractors.ID(r, "id")
  if err != nil {
    return err
  }

  entity, err := uc.service.Get(r.Context(), id)
  if err != nil {
    return err
  }

  return server.JSON(w, http.StatusOK, entity)
}


// Create godoc
//
// @Tags      {{ .Computed.domain_var }}
// @Summary   Create a new {{ .Computed.domain_var }}
// @Description Create a new {{ .Computed.domain_var }}
// @Accept     json
// @Produce    json
// @Param      {{ .Computed.domain_var }} body    dtos.{{ .Computed.domain_var }}Create  true  "The {{ .Computed.domain_var }} details"
// @Success    201     {object}  dtos.{{ .Computed.domain_var }}
// @Router     /api/v1/{{ .Computed.domain_kebab }}s [POST]
func (uc *{{ .Computed.domain_var }}Controller) Create(w http.ResponseWriter, r *http.Request) error {
  body, err := extractors.Body[dtos.{{ .Computed.domain_var }}Create](r)
  if err != nil {
    return err
  }

  entity, err := uc.service.Create(r.Context(), body)
  if err != nil {
    return err
  }

  return server.JSON(w, http.StatusCreated, entity)
}

// Update godoc
//
// @Tags      {{ .Computed.domain_var }}
// @Summary   Update a {{ .Computed.domain_var }}
// @Description Update a {{ .Computed.domain_var }}
// @Accept     json
// @Produce    json
// @Param      {{ .Computed.domain_var }} body    dtos.{{ .Computed.domain_var }}Update  true  "The {{ .Computed.domain_var }} details"
// @Param      id      path    string  true  "The {{ .Computed.domain_var }} ID"
// @Success    200     {object}  dtos.{{ .Computed.domain_var }}
// @Router     /api/v1/{{ .Computed.domain_kebab }}s/{id} [PUT]
func (uc *{{ .Computed.domain_var }}Controller) Update(w http.ResponseWriter, r *http.Request) error {
  id, body, err := extractors.BodyWithID[dtos.{{ .Computed.domain_var }}Update](r, "id")
  if err != nil {
    return err
  }

  entity, err := uc.service.Update(r.Context(), id, body)
  if err != nil {
    return err
  }

  return server.JSON(w, http.StatusOK, entity)
}

// Delete godoc
//
// @Tags      {{ .Computed.domain_var }}
// @Summary   Delete a {{ .Computed.domain_var }}
// @Description Delete a {{ .Computed.domain_var }}
// @Accept     json
// @Produce    json
// @Param      id      path    string  true  "The {{ .Computed.domain_var }} ID"
// @Success    204     
// @Router     /api/v1/{{ .Computed.domain_kebab }}s/{id} [DELETE]
func (uc *{{ .Computed.domain_var }}Controller) Delete(w http.ResponseWriter, r *http.Request) error {
  id, err := extractors.ID(r, "id")
  if err != nil {
    return err
  }

  err = uc.service.Delete(r.Context(), id)
  if err != nil {
    return err
  }

  return server.JSON(w, http.StatusNoContent, nil)
}

