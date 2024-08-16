package main

import (
	"context"
	"time"

	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	Email    string `conf:"default:admin@example.com"`
	Username string `conf:"default:admin"`
	Password string `conf:"default:admin1"`
}

func seed(q *db.QueriesExt, cfg Seed) error {
	svcs := services.NewService(services.Config{}, log.Logger, q, nil)

	user, err := svcs.Admin.Register(context.Background(), dtos.UserRegister{
		Email:    cfg.Email,
		Username: cfg.Username,
		Password: cfg.Password,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to seed database")
		return err
	}

	var (
		subCustID    = "cus_ADMIN_OVERRIDE"
		subID        = "sub_ADMIN_OVERRIDE"
		subStartDate = time.Now()
		subEndDate   = time.Now().AddDate(20, 0, 0)
	)

	_, err = svcs.Users.UpdateSubscription(context.Background(), user.ID, dtos.UserUpdateSubscription{
		StripeCustomerID:      &subCustID,
		StripeSubscriptionID:  &subID,
		SubscriptionStartDate: &subStartDate,
		SubscriptionEndedDate: &subEndDate,
	})
	if err != nil {
		log.Error().Err(err).Str("step", "subscription data").Msg("failed to seed database")
		return err
	}

	log.Info().
		Str("id", user.ID.String()).
		Str("email", user.Email).
		Str("username", user.Username).
		Msg("successfully seeded database")

	return nil
}
