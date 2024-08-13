package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
)

type UserService struct {
	l      zerolog.Logger
	db     *db.QueriesExt
	mapper dtos.MapFunc[db.User, dtos.User]
}

func NewUserService(l zerolog.Logger, db *db.QueriesExt) *UserService {
	return &UserService{
		l:      l,
		db:     db,
		mapper: dtos.MapUser,
	}
}

func (s *UserService) Register(ctx context.Context, data dtos.UserRegister) (dtos.User, error) {
	panic("TODO")
}

// Get returns a single user by id from the database.
func (s *UserService) Get(ctx context.Context, id uuid.UUID) (dtos.User, error) {
	v, err := s.db.UserByID(ctx, id)
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}

// Delete _fully_ deletes a user from the database and all associated content. This
// should be considered full account deletion with ABSOLUTELY NO WAY to undo these
// changes.
//
// Use with caution
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.UserDeleteByID(ctx, id)
}

// UpdateDetails updates the user's username and email address. If the value is null, it will
// be ignored during the setting within the database.
func (s *UserService) UpdateDetails(ctx context.Context, id uuid.UUID, data dtos.UserUpdate) (dtos.User, error) {
	v, err := s.db.UserUpdate(ctx, db.UserUpdateParams{
		ID:           id,
		Username:     data.Usenrame,
		Email:        data.Email,
		PasswordHash: nil,
	})
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}

// UpdateSubscription updates the user's subscription details
func (s *UserService) UpdateSubscription(ctx context.Context, id uuid.UUID, data dtos.UserUpdateSubscription) (dtos.User, error) {
	v, err := s.db.UserUpdateBilling(ctx, db.UserUpdateBillingParams{
		ID:                    id,
		StripeCustomerID:      data.StripeCustomerID,
		StripeSubscriptionID:  data.StripeSubscriptionID,
		SubscriptionStartDate: db.IntoPgTimePrt(data.SubscriptionStartDate),
		SubscriptionEndedDate: db.IntoPgTimePrt(data.SubscriptionEndedDate),
	})
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}
