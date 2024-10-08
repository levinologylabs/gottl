package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jalevin/gottl/internal/core/hasher"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/rs/zerolog"
)

var ErrInvalidLogin = errors.New("invalid password or username")

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
	pwHash, err := hasher.HashPassword(data.Password)
	if err != nil {
		return dtos.User{}, err
	}

	v, err := s.db.UserCreate(ctx, db.UserCreateParams{
		Username:     data.Username,
		Email:        data.Email,
		PasswordHash: pwHash,
	})
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}

// Authenticate validates a user's credentials and returns the user if they are valid.
// If the credentials are invalid, an error is returned. This function uses a constant
// time comparison to prevent timing attacks. When no use is found by the provided email
// address, the same error is returned to prevent user enumeration.
func (s *UserService) Authenticate(ctx context.Context, data dtos.UserAuthenticate) (dtos.UserSession, error) {
	dbuser, err := s.db.UserByEmail(ctx, data.Email)
	if err != nil {
		// This is to prevent timing attacks ensuring that when no user is found we
		// still perform the same amount of work as when a user is found.
		//
		// savedHash = ThisIsNotAStrongPassword?12!#$%@!@!@$ButItWorks
		savedHash := "$argon2id$v=19$m=65536,t=1,p=8$r14KLB8NUVfFFccYbU1q9w$tJ3HvNwMED2dL3lALmOdkm46TVuB9vGcEjy9sxTAE6s"
		hasher.CheckPasswordHash(data.Password, savedHash)

		s.l.Error().Ctx(ctx).Err(err).Str("email", data.Email).Msg("failed to get user by email")
		return dtos.UserSession{}, ErrInvalidLogin
	}

	if !hasher.CheckPasswordHash(data.Password, dbuser.PasswordHash) {
		s.l.Error().Ctx(ctx).Err(err).Str("email", data.Email).Msg("password verification failed")
		return dtos.UserSession{}, ErrInvalidLogin
	}

	return s.createSession(ctx, s.mapper.Map(dbuser))
}

// SessionVerify validates a user's session token and returns the user if the token is valid
// and has not expired.
func (s *UserService) SessionVerify(ctx context.Context, token string) (dtos.User, error) {
	user, err := s.db.UserBySession(ctx, hasher.HashToken(token))
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(user), nil
}

func (s *UserService) createSession(ctx context.Context, user dtos.User) (dtos.UserSession, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 31)

	token := hasher.NewToken()

	err := s.db.SessionCreate(ctx, db.SessionCreateParams{
		UserID:    user.ID,
		Token:     token.Hash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return dtos.UserSession{}, err
	}

	return dtos.UserSession{
		Token:     token.Raw,
		ExpiresAt: expiresAt,
	}, nil
}

// GetByID returns a single user by id from the database.
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (dtos.User, error) {
	v, err := s.db.UserByID(ctx, id)
	if err != nil {
		return dtos.User{}, err
	}

	return s.mapper.Map(v), nil
}

// GetByEmail returns a single user by email from the database.
func (s *UserService) GetByEmail(ctx context.Context, email string) (dtos.User, error) {
	v, err := s.db.UserByEmail(ctx, email)
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
		Username:     data.Username,
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

const OAuthPasswordPlaceholder = "ProviderOnlyUser"

// ProviderSession creates a new session for a user that has authenticated with a third-party.
func (s *UserService) ProviderSession(
	ctx context.Context,
	providerName string,
	extID string,
	extEmail string,
	extName string,
) (dtos.UserSession, error) {
	var (
		err    error
		dbuser db.User
		user   dtos.User
	)

	// try get user by extID
	dbuser, err = s.db.UserByProvider(ctx, db.UserByProviderParams{
		ProviderName:   providerName,
		ProviderUserID: extID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return dtos.UserSession{}, err
	}

	if dbuser.ID == uuid.Nil {
		// try bind by email
		dbuser, err = s.db.UserByEmail(ctx, extEmail)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return dtos.UserSession{}, err
		}
	}

	if dbuser.ID != uuid.Nil {
		user = s.mapper.Map(dbuser)
	}

	// If the user is still a zero value at this point, we need to create a new user
	// and bind the provider to that user.
	if user.ID == uuid.Nil {
		user, err = s.Register(ctx, dtos.UserRegister{
			Email:    extEmail,
			Username: uuid.NewString(), // Randomly generate a username since we can't guarantee one from the provider is unique
			Password: OAuthPasswordPlaceholder,
		})
		if err != nil {
			s.l.Err(err).Msg("error creating user during oauth")
			return dtos.UserSession{}, err
		}

		_, err = s.db.CreateProvider(ctx, db.CreateProviderParams{
			UserID:         user.ID,
			ProviderName:   providerName,
			ProviderUserID: extID,
			Metadata:       nil,
		})
		if err != nil {
			s.l.Err(err).Msg("error creating provider during oauth")
			return dtos.UserSession{}, err
		}
	}

	if user.ID == uuid.Nil {
		panic("unable to resolve user during oauth validation")
	}

	return s.createSession(ctx, user)
}

func (s *UserService) ProviderStateGet(ctx context.Context) (string, error) {
	token := hasher.NewToken()

	err := s.db.ProviderStateCreate(ctx, db.ProviderStateCreateParams{
		Token:     token.Hash,
		ExpiresAt: time.Now().Add(time.Minute * 5),
	})
	if err != nil {
		return "", err
	}

	return token.Raw, nil
}

func (s *UserService) ProviderStateUse(ctx context.Context, token string) error {
	hash := hasher.HashToken(token)
	return s.db.WithTx(ctx, func(qe *db.QueriesExt) error {
		_, err := qe.ProviderStateGet(ctx, hash)
		if err != nil {
			return err
		}

		err = qe.ProviderStateDelete(ctx, hash)
		if err != nil {
			return err
		}

		return nil
	})
}
