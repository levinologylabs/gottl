package services_test

import (
	"context"
	"testing"

	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PasswordService_ResetPassword(t *testing.T) {
	testlib.IntegrationGuard(t)

	tctx := SetupServiceTest(t)
	svc := services.NewPasswordService(services.Config{}, tctx.logger, tctx.db, tasks.NoopQueue)

	token, err := svc.RequestReset(context.Background(), dtos.PasswordResetRequest{
		Email: tctx.user.Email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	for i := range 2 {
		err = svc.Reset(context.Background(), dtos.PasswordReset{
			Token:    token,
			Password: "newpassword",
		})

		switch i {
		case 0:
			require.NoError(t, err)
		case 1:
			require.Error(t, err, "expected token to be used on second try")
		}
	}

	userservice := services.NewUserService(tctx.logger, tctx.db)
	_, err = userservice.Authenticate(context.Background(), dtos.UserAuthenticate{
		Email:    tctx.user.Email,
		Password: "newpassword",
	})
	require.NoError(t, err)
}
