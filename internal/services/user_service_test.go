package services_test

import (
	"context"
	"testing"

	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserService_RegisterAndLogin(t *testing.T) {
	testlib.IntegrationGuard(t)

	logger := testlib.TestLogger(t)
	queries := testlib.NewDatabase(t, logger)

	s := services.NewUserService(logger, queries)
	ctx := context.Background()

	registerUser, err := s.Register(ctx, dtos.UserRegister{
		Email:    "user@example.com",
		Username: "user",
		Password: "MyUserPassword?12",
	})

	require.NoError(t, err)
	require.NotNil(t, registerUser)

	loginUser, err := s.Authenticate(ctx, dtos.UserAuthenticate{
		Email:    "user@example.com",
		Password: "MyUserPassword?12",
	})

	require.NoError(t, err)
	require.NotNil(t, loginUser)

	assert.Equal(t, registerUser, loginUser)
}
