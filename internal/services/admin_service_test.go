package services_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type AdminServiceTestCtx struct {
	us        UserServiceTestCtx
	ctx       context.Context
	s         *services.AdminService
	adminuser dtos.UserRegister
}

func SetupAdminServiceTest(t *testing.T) AdminServiceTestCtx {
	us := SetupUserServiceTest(t)
	var (
		logger  = testlib.Logger(t)
		queries = testlib.NewDatabase(t, logger)
		s       = services.NewAdminService(logger, queries)
	)
	return AdminServiceTestCtx{
		us:  us,
		ctx: context.Background(),
		s:   s,
		adminuser: dtos.UserRegister{
			Email:    gofakeit.Email(),
			Username: gofakeit.Username(),
			Password: gofakeit.Password(true, true, true, true, true, 14),
		},
	}
}

func Test_AdminService_GetAllUsers(t *testing.T) {
	testlib.IntegrationGuard(t)
	st := SetupAdminServiceTest(t)

	// Regular User cannot access method
	_, err := st.s.GetAllUsers(st.ctx, dtos.Pagination{}.WithDefaults())
	require.Error(t, err)

	_, err = st.s.Register(st.ctx, st.adminuser)
	require.NoError(t, err)

	// Admin User can access method
	adminctx := services.WithVerifiedAdmin(st.ctx)
	users, err := st.s.GetAllUsers(adminctx, dtos.Pagination{}.WithDefaults())
	require.NoError(t, err)

	assert.Len(t, users.Items, 1)
	assert.Equal(t, len(users.Items), users.Total)
}
