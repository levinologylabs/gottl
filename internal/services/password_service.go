package services

import (
	"context"
	"time"

	"github.com/jalevin/gottl/internal/core/hasher"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
)

type PasswordService struct {
	l  zerolog.Logger
	db *db.QueriesExt
}

func NewPasswordService(l zerolog.Logger, db *db.QueriesExt) *PasswordService {
	return &PasswordService{
		l:  l,
		db: db,
	}
}

func (s *PasswordService) RequestReset(ctx context.Context, data dtos.PasswordResetRequest) error {
	// launch background task to send email
	panic("TODO")
}

func (s *PasswordService) Reset(ctx context.Context, data dtos.PasswordReset) error {
	return s.db.WithTx(ctx, func(qe *db.QueriesExt) error {
		v, err := s.db.UserActionTokenGet(ctx, db.UserActionTokenGetParams{
			Token:  hasher.HashToken(data.Token),
			Action: db.UserActionPasswordReset,
			Now:    time.Now(),
		})
		if err != nil {
			return err
		}

		pwHash, err := hasher.HashPassword(data.Password)
		if err != nil {
			return err
		}

		_, err = s.db.UserUpdate(ctx, db.UserUpdateParams{
			PasswordHash: &pwHash,
		})
		if err != nil {
			return err
		}

		err = s.db.UserActionTokenDelete(ctx, v.ID)
		if err != nil {
			return err
		}

		return nil
	})
}
