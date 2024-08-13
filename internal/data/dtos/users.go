package dtos

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`

	StripeCustomerID      *string   `json:"-"`
	StripeSubscriptionID  *string   `json:"-"`
	SubscriptionStartDate time.Time `json:"-"`
	SubscriptionEndedDate time.Time `json:"-"`
}

// IsSubscribed is a proxy for IsSubscribedAt(time.Now())
func (u User) IsSubscribed() bool {
	return u.IsSubscribedAt(time.Now())
}

// IsSubscribedAt returns true if the user is subscribed at the given time
// provided. This is based on the subscription date range and does not validate
// the underlying stripe details.
//
// Use this to determine if a user is subscribed, as it will allow you to bypass
// the stripe requirement for test or admin accounts.
func (u User) IsSubscribedAt(now time.Time) bool {
	if u.SubscriptionStartDate.Before(now) && u.SubscriptionEndedDate.After(now) {
		return true
	}

	return false
}

type UserRegister struct {
	Email    string `json:"email"    validate:"required,email"`
	Username string `json:"username" validate:"required,min=6,max=128"`
	Password string `json:"password" validate:"required,min=6,max=256"`
}

type UserUpdate struct {
	Email    *string `json:"email"    validate:"omitempty,email"`
	Usenrame *string `json:"username" validate:"omitempty,min=6,max=128"`
	Password *string `json:"password" validate:"omitempty,min=6,max=256"`
}

type UserUpdateSubscription struct {
	StripeCustomerID      *string    `json:"stripeCustomerId"`
	StripeSubscriptionID  *string    `json:"stripeSubscriptionId"`
	SubscriptionStartDate *time.Time `json:"-"`
	SubscriptionEndedDate *time.Time `json:"-"`
}
