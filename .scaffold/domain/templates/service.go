package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
)

type {{ .Computed.domain_var }}Service struct {
  l zerolog.Logger
  db *db.QueriesExt
  mapper dtos.MapFunc[db.{{ .Computed.domain_var }}, dtos.{{ .Computed.domain_var }}]
}

func New{{ .Computed.domain_var }}Service(l zerolog.Logger, db *db.QueriesExt) *{{ .Computed.domain_var }}Service {
  return &{{ .Computed.domain_var }}Service{
    l: l,
    db: db,
    mapper: dtos.Map{{ .Computed.domain_var }},
  }
}

func (s *{{ .Computed.domain_var}}Service) Get(ctx context.Context, id uuid.UUID) (dtos.{{ .Computed.domain_var }}, error) {
  entity, err := s.db.{{ .Computed.domain_var }}ByID(ctx, id)
  if err != nil {
    return dtos.{{ .Computed.domain_var }}{}, err
  }

  return s.mapper(entity), nil
}

func (s *{{ .Computed.domain_var}}Service) GetAll(ctx context.Context, page dtos.Pagination) (dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}], error) {
  count, err := s.db.{{ .Computed.domain_var }}GetAllCount(ctx)
  if err != nil {
    return dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}]{}, err
  }

  entities, err := s.db.{{ .Computed.domain_var }}GetAll(ctx, db.{{ .Computed.domain_var }}GetAllParams{
    Limit: int32(page.Limit),
    Offset: int32(page.Skip),
  })
  if err != nil {
    return dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}]{}, err
  }

  return dtos.PaginationResponse[dtos.{{ .Computed.domain_var }}]{
    Total: int(count),
    Items: s.mapper.Slice(entities),
  }, nil
}

func (s *{{ .Computed.domain_var}}Service) Create(ctx context.Context, data dtos.{{ .Computed.domain_var }}Create) (dtos.{{ .Computed.domain_var }}, error) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) Update(ctx context.Context, id uuid.UUID, data dtos.{{ .Computed.domain_var }}Update) (dtos.{{ .Computed.domain_var }}, error) {
  panic("not implemented")
}

func (s *{{ .Computed.domain_var}}Service) Delete(ctx context.Context, id uuid.UUID) error {
  err := s.db.{{ .Computed.domain_var }}DeleteByID(ctx, id)
  if err != nil {
    return err
  }

  return nil
}
