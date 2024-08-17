package dtos

import (
	"time"

	"github.com/google/uuid"
)

type {{ .Computed.domain_var }} struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type {{ .Computed.domain_var }}Create struct {}

type {{ .Computed.domain_var }}Update struct {}
