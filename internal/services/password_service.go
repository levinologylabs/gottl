package services

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jalevin/gottl/internal/core/hasher"
	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services/emailtemplates"
	"github.com/rs/zerolog"
)

type PasswordService struct {
	cfg   Config
	l     zerolog.Logger
	db    *db.QueriesExt
	queue tasks.Queue
}

func NewPasswordService(cfg Config, l zerolog.Logger, db *db.QueriesExt, queue tasks.Queue) *PasswordService {
	return &PasswordService{
		l:     l,
		db:    db,
		queue: queue,
	}
}

func (s *PasswordService) RequestReset(ctx context.Context, data dtos.PasswordResetRequest) error {
	usr, err := s.db.UserByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.l.Warn().Str("email", data.Email).Msg("user not found for password reset request")
			return nil
		}
	}

	token := hasher.NewToken()
	_, err = s.db.UserActionTokenCreate(ctx, db.UserActionTokenCreateParams{
		UserID:    usr.ID,
		Token:     token.Hash,
		Action:    db.UserActionPasswordReset,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}

	return s.queue.Enqueue(tasks.Task{
		ID: tasks.TaskIDSendEmail,
		Payload: tasks.TaskDataSendEmail{
			Email:   data.Email,
			Subject: "Password Reset",
			Body:    emailtemplates.PasswordReset(s.cfg.CompanyName, s.cfg.WebURL, token.Raw),
		},
	})
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
