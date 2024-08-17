package services_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/jalevin/gottl/internal/data/dtos"
	"github.com/jalevin/gottl/internal/services"
	"github.com/jalevin/gottl/testlib"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type ServiceTestCtx struct {
	db      *db.QueriesExt
	logger  zerolog.Logger
	user    dtos.UserRegister
	dbuser  dtos.User
	admin   dtos.UserRegister
	dbadmin dtos.User
}

func SetupServiceTest(t *testing.T) ServiceTestCtx {
	t.Helper()

	var (
		logger  = testlib.Logger(t)
		queries = testlib.NewDatabase(t, logger)
	)

	svcuser := services.NewUserService(logger, queries)
	svcadmin := services.NewAdminService(logger, queries)

	user := dtos.UserRegister{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 14),
	}

	admin := dtos.UserRegister{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 14),
	}

	dbuser, err := svcuser.Register(context.Background(), user)
	require.NoError(t, err)

  dbadmin, err := svcadmin.Register(context.Background(), admin)
	require.NoError(t, err)

	return ServiceTestCtx{
		db:      queries,
		logger:  logger,
		user:    user,
		admin:   admin,
		dbuser:  dbuser,
		dbadmin: dbadmin,
	}
}
