package user_test

import (
	"database/sql"
	"testing"
	"xm/configs"
	"xm/pkg/logger"
	userRepo "xm/pkg/repositories/user"
	"xm/pkg/services"
	"xm/pkg/services/user"
	"xm/pkg/services/utils"

	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreate(t *testing.T) {
	svc, m := getTestService(t)

	req := userRepo.User{
		Username: "some",
		Password: "pass",
	}

	m.On("Create", &req).Return(nil).Once()

	err := svc.Create(&req)
	require.NoError(t, err)

	m.On("Create", &req).Return(error(&pq.Error{Code: "23505"}))
	err = svc.Create(&req)
	require.ErrorIs(t, err, utils.ErrAlreadyExists)
}

func TestGetByUsername(t *testing.T) {
	svc, m := getTestService(t)

	req := userRepo.User{
		Username: "some",
		Password: "pass",
	}

	m.On("GetByUsername", req.Username).Return(&req, nil)

	u, err := svc.GetByUsername(req.Username)
	require.NoError(t, err)

	require.Equal(t, u.Password, req.Password)

	req = userRepo.User{
		Username: "not found",
	}

	m.On("GetByUsername", req.Username).Return(nil, sql.ErrNoRows)

	u, err = svc.GetByUsername(req.Username)
	require.Equal(t, err, utils.ErrNotFound)

	require.Nil(t, u)
}

func getTestService(t *testing.T) (user.Service, *mocker) {
	var repo user.Service
	m := &mocker{}

	go fxtest.New(
		fxtest.TB(t),
		fx.Options(
			configs.Module,
			logger.Module,

			fx.Provide(
				func() userRepo.Repository {
					return m
				},
			),

			services.Module,
		),
		fx.Populate(&repo),
	).Run()

	return repo, m
}

type mocker struct {
	mock.Mock
}

func (m *mocker) Create(u *userRepo.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mocker) GetByUsername(username string) (*userRepo.User, error) {
	var u *userRepo.User

	args := m.Called(username)
	u, _ = args.Get(0).(*userRepo.User)

	return u, args.Error(1)
}
