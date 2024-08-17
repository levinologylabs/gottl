package services


import (
  "context"

	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
  "github.com/rs/zerolog"
)

type {{ .Computed.domain_var }}Service struct {
  l zerolog.Logger
  db *db.QueriesExt
}

func New{{ .Computed.domain_var }}Service(l zerolog.Logger, db *db.QueriesExt) *{{ .Computed.domain_var }}Service {
  return &{{ .Computed.domain_var }}Service{
    l: l,
    db: db,
  }
}

func (s *{{ .Computed.domain_var}}Service) GetOne(ctx context.Contet, id uuid.UUID) (dtos.{{ .Computed.domain_var }}) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) GetAll(ctx context.Contet, page dtos.Pagination) (dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}]) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) Create(ctx context.Contet, data dtos.{{ .Computed.domain_var }}Create) (dtos.{{ .Computed.domain_var }}) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) Update(ctx context.Contet, id uuid.UUID, data dtos.{{ .Computed.domain_var }}Update) (dtos.{{ .Computed.domain_var }}) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) Delete(ctx context.Contet, id uuid.UUID) {
  panic("not implemented")
}
