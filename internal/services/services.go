// Package services contains the main business logic of the application
package services

import (
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/rs/zerolog"
)

// Service is a collection of all services in the application
type Service struct {
	Admin *AdminService
	Users *UserService
}

func NewService(l zerolog.Logger, db *db.QueriesExt) *Service {
	return &Service{
		Admin: NewAdminService(l, db),
		Users: NewUserService(l, db),
	}
}
