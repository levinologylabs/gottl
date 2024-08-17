package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserService_RegisterAndLogin(t *testing.T) {
	testlib.IntegrationGuard(t)
	sst := SetupServiceTest(t)
	st := services.NewUserService(sst.logger, sst.db)

	user := dtos.UserRegister{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 14),
	}

	registerUser, err := st.Register(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, user.Email, registerUser.Email)
	assert.Equal(t, user.Username, registerUser.Username)

	// Login With Correct Password

	session, err := st.Authenticate(context.Background(), dtos.UserAuthenticate{
		Email:    user.Email,
		Password: user.Password,
	})

	require.NoError(t, err)
	require.NotNil(t, session)

	assert.NotEmpty(t, session.Token)
	assert.True(t, session.ExpiresAt.After(time.Now()))

	loginUser, err := st.SessionVerify(context.Background(), session.Token)
	require.NoError(t, err)

	assert.Equal(t, registerUser, loginUser)

	// Login With Wrong Password Fails
	_, err = st.Authenticate(context.Background(), dtos.UserAuthenticate{
		Email:    user.Email,
		Password: "wrongpassword",
	})

	require.Error(t, err)
}

func Test_UserService_GetUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		sst = SetupServiceTest(t)
		st  = services.NewUserService(sst.logger, sst.db)
	)

	userByEmail, err := st.GetByEmail(context.Background(), sst.user.Email)
	require.NoError(t, err)

	userByID, err := st.GetByID(context.Background(), userByEmail.ID)
	require.NoError(t, err)

	assert.Equal(t, userByEmail, userByID)
}

func Test_UserService_DeleteUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		st   = SetupServiceTest(t)
		s    = services.NewUserService(st.logger, st.db)
		user = st.dbuser
	)

	err := s.Delete(context.Background(), user.ID)
	require.NoError(t, err)

	_, err = s.GetByID(context.Background(), user.ID)
	require.Error(t, err)
}

func Test_UserService_UpdateUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		st   = SetupServiceTest(t)
		s    = services.NewUserService(st.logger, st.db)
		user = st.dbuser
	)

	// Patch the user's email address
	updatedUser, err := s.UpdateDetails(context.Background(), user.ID, dtos.UserUpdate{
		Email:    testlib.Ptr("new@example.com"),
		Username: nil,
		Password: nil,
	})
	require.NoError(t, err)

	assert.Equal(t, "new@example.com", updatedUser.Email, "email should have changed")
	assert.Equal(t, user.Username, updatedUser.Username, "username should not have changed")

	// Patch Subscription Details
	subdata := dtos.UserUpdateSubscription{
		StripeCustomerID:      testlib.Ptr("cus_12345"),
		StripeSubscriptionID:  testlib.Ptr("sub_12345"),
		SubscriptionStartDate: testlib.Ptr(time.Now().UTC()),
		SubscriptionEndedDate: testlib.Ptr(time.Now().Add(time.Hour * 24 * 30).UTC()),
	}

	updatedUser, err = s.UpdateSubscription(context.Background(), user.ID, subdata)

	require.NoError(t, err)

	assert.Equal(t, subdata.StripeCustomerID, updatedUser.StripeCustomerID)
	assert.Equal(t, subdata.StripeSubscriptionID, updatedUser.StripeSubscriptionID)
	assert.Equal(t, subdata.SubscriptionStartDate.UnixMilli(), updatedUser.SubscriptionStartDate.UnixMilli())
	assert.Equal(t, subdata.SubscriptionEndedDate.UnixMilli(), updatedUser.SubscriptionEndedDate.UnixMilli())
}
