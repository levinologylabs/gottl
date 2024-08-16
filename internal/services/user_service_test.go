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

// UserServiceTestCtx is a container type for housing common test data and setup login
// for the user service tests.
type UserServiceTestCtx struct {
	ctx  context.Context
	s    *services.UserService
	user dtos.UserRegister
}

func SetupUserServiceTest(t *testing.T) UserServiceTestCtx {
	var (
		base = SetupServiceTest(t)
		s    = services.NewUserService(base.logger, base.db)
	)

	return UserServiceTestCtx{
		ctx: context.Background(),
		s:   s,
		user: dtos.UserRegister{
			Email:    gofakeit.Email(),
			Username: gofakeit.Username(),
			Password: gofakeit.Password(true, true, true, true, true, 14),
		},
	}
}

func (ts *UserServiceTestCtx) RegisterUser(t *testing.T) dtos.User {
	registerUser, err := ts.s.Register(context.Background(), ts.user)
	require.NoError(t, err)
	return registerUser
}

func Test_UserService_RegisterAndLogin(t *testing.T) {
	testlib.IntegrationGuard(t)
	st := SetupUserServiceTest(t)

	registerUser := st.RegisterUser(t)
	assert.Equal(t, st.user.Email, registerUser.Email)
	assert.Equal(t, st.user.Username, registerUser.Username)

	// Login With Correct Password

	session, err := st.s.Authenticate(context.Background(), dtos.UserAuthenticate{
		Email:    st.user.Email,
		Password: st.user.Password,
	})

	require.NoError(t, err)
	require.NotNil(t, session)

	assert.NotEmpty(t, session.Token)
	assert.True(t, session.ExpiresAt.After(time.Now()))

	loginUser, err := st.s.SessionVerify(st.ctx, session.Token)
	require.NoError(t, err)

	assert.Equal(t, registerUser, loginUser)

	// Login With Wrong Password Fails
	_, err = st.s.Authenticate(context.Background(), dtos.UserAuthenticate{
		Email:    st.user.Email,
		Password: "wrongpassword",
	})

	require.Error(t, err)
}

func Test_UserService_GetUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		st   = SetupUserServiceTest(t)
		user = st.RegisterUser(t)
	)

	userByID, err := st.s.GetByID(st.ctx, user.ID)
	require.NoError(t, err)

	assert.Equal(t, user, userByID)

	userByEmail, err := st.s.GetByEmail(st.ctx, user.Email)
	require.NoError(t, err)

	assert.Equal(t, user, userByEmail)
}

func Test_UserService_DeleteUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		st   = SetupUserServiceTest(t)
		user = st.RegisterUser(t)
	)

	err := st.s.Delete(st.ctx, user.ID)
	require.NoError(t, err)

	_, err = st.s.GetByID(st.ctx, user.ID)
	require.Error(t, err)
}

func Test_UserService_UpdateUser(t *testing.T) {
	testlib.IntegrationGuard(t)

	var (
		st   = SetupUserServiceTest(t)
		user = st.RegisterUser(t)
	)

	// Patch the user's email address
	updatedUser, err := st.s.UpdateDetails(st.ctx, user.ID, dtos.UserUpdate{
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

	updatedUser, err = st.s.UpdateSubscription(st.ctx, user.ID, subdata)

	require.NoError(t, err)

	assert.Equal(t, subdata.StripeCustomerID, updatedUser.StripeCustomerID)
	assert.Equal(t, subdata.StripeSubscriptionID, updatedUser.StripeSubscriptionID)
	assert.Equal(t, subdata.SubscriptionStartDate.UnixMilli(), updatedUser.SubscriptionStartDate.UnixMilli())
	assert.Equal(t, subdata.SubscriptionEndedDate.UnixMilli(), updatedUser.SubscriptionEndedDate.UnixMilli())
}
