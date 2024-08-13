package dtos

import "github.com/jalevin/gottl/internal/data/db"

// MapUser maps a db user into a dto user type
func MapUser(dbu db.User) User {
	return User{
		ID:                    dbu.ID,
		CreatedAt:             dbu.CreatedAt,
		UpdatedAt:             dbu.UpdatedAt,
		Username:              dbu.Username,
		Email:                 dbu.Email,
		StripeCustomerID:      dbu.StripeCustomerID,
		StripeSubscriptionID:  dbu.StripeSubscriptionID,
		SubscriptionStartDate: dbu.SubscriptionStartDate.Time,
		SubscriptionEndedDate: dbu.SubscriptionEndedDate.Time,
	}
}
