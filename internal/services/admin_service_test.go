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

func Test_AdminService_GetAllUsers(t *testing.T) {
	testlib.IntegrationGuard(t)
	var (
		st = SetupServiceTest(t)
		s  = services.NewAdminService(st.logger, st.db)
	)

	// Regular User cannot access method
	_, err := s.GetAllUsers(context.Background(), dtos.Pagination{}.WithDefaults())
	require.Error(t, err)

	// Admin User can access method
	adminctx := services.WithVerifiedAdmin(context.Background())
	users, err := s.GetAllUsers(adminctx, dtos.Pagination{}.WithDefaults())
	require.NoError(t, err)

	// Two users created in SetupServiceTest
	assert.Len(t, users.Items, 2)
}
