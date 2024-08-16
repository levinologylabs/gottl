package services

import (
	"context"

	"github.com/jalevin/gottl/internal/core/hasher"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
)

type AdminService struct {
	l      zerolog.Logger
	db     *db.QueriesExt
	mapper dtos.MapFunc[db.User, dtos.User]
}

func NewAdminService(l zerolog.Logger, db *db.QueriesExt) *AdminService {
	return &AdminService{
		l:      l,
		db:     db,
		mapper: dtos.MapUser,
	}
}

func (s *AdminService) Register(ctx context.Context, date dtos.UserRegister) (dtos.User, error) {
	pwHash, err := hasher.HashPassword(date.Password)
	if err != nil {
		return dtos.User{}, err
	}

	v, err := s.db.UserCreateAdmin(ctx, db.UserCreateAdminParams{
		Username:     date.Username,
		Email:        date.Email,
		PasswordHash: pwHash,
	})
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}

func (s *AdminService) GetAllUsers(ctx context.Context, page dtos.Pagination) (dtos.PaginationResponse[dtos.User], error) {
	var out dtos.PaginationResponse[dtos.User]
	if err := VerifiedAdminFrom(ctx); err != nil {
		return out, err
	}

	err := s.db.WithTx(ctx, func(q *db.QueriesExt) error {
		v, err := s.db.UserGetAll(ctx, db.UserGetAllParams{
			Limit:  int32(page.Limit),
			Offset: int32(page.Skip),
		})
		if err != nil {
			return err
		}

		count, err := s.db.UserGetAllCount(ctx)
		if err != nil {
			return err
		}

		out = dtos.PaginationResponse[dtos.User]{
			Items: s.mapper.Slice(v),
			Total: int(count),
		}

		return nil
	})
	if err != nil {
		return out, err
	}

	return out, nil
}
