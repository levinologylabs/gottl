// Package services contains the main business logic of the application
package services

import (
	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/rs/zerolog"
)

// Service is a collection of all services in the application
type Service struct {
	Admin     *AdminService
	Users     *UserService
	Passwords *PasswordService
  // $scaffold_inject_service
}

func NewService(
	cfg Config,
	l zerolog.Logger,
	db *db.QueriesExt,
	queue tasks.Queue,
) *Service {
	return &Service{
		Admin:     NewAdminService(l, db),
		Users:     NewUserService(l, db),
		Passwords: NewPasswordService(cfg, l, db, queue),
    // $scaffold_inject_service_constructor
	}
}
