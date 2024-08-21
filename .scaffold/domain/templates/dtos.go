package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/jalevin/gottl/internal/data/db"
)

type {{ .Computed.domain_var }} struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type {{ .Computed.domain_var }}Create struct {}

type {{ .Computed.domain_var }}Update struct {}

func Map{{ .Computed.domain_var }}(d db.{{ .Computed.domain_var }}) {{ .Computed.domain_var }} {
  return {{ .Computed.domain_var }}{
    ID:        d.ID,
    CreatedAt: d.CreatedAt,
    UpdatedAt: d.UpdatedAt,
  }
}
