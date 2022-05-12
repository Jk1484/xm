package company_test

import (
	"database/sql"
	"testing"
	"xm/configs"
	"xm/gateways"
	"xm/pkg/logger"
	"xm/pkg/services"
	"xm/pkg/services/company"
	"xm/pkg/services/utils"

	companyRepo "xm/pkg/repositories/company"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreate(t *testing.T) {
	svc, m := getTestService(t)

	cmp := companyRepo.Company{
		Name:    "name",
		Code:    "code",
		Country: "country",
	}

	m.On("Create", cmp).Return(nil)

	err := svc.Create(cmp)
	require.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	svc, m := getTestService(t)

	m.On("GetByID", 1).Return(companyRepo.Company{Name: "some company"}, nil)

	c, err := svc.GetByID(1)
	require.NoError(t, err)
	require.Equal(t, c.Name, "some company")

	m.On("GetByID", 2).Return(companyRepo.Company{}, sql.ErrNoRows)

	c, err = svc.GetByID(2)
	require.ErrorIs(t, err, utils.ErrNotFound)
	require.Zero(t, c)
}

func TestGetAll(t *testing.T) {
	svc, m := getTestService(t)

	f := companyRepo.Filters{}

	m.On("GetAll", f).Return([]companyRepo.Company{{}, {}}, nil).Once()

	cs, err := svc.GetAll(f)
	require.NoError(t, err)
	require.Equal(t, len(cs), 2)

	m.On("GetAll", f).Return(nil, sql.ErrNoRows)

	cs, err = svc.GetAll(f)
	require.ErrorIs(t, err, utils.ErrNotFound)
	require.Nil(t, cs)
}

func TestUpdate(t *testing.T) {
	svc, m := getTestService(t)

	c := companyRepo.Company{}

	m.On("Update", c).Return(nil)

	err := svc.Update(c)
	require.NoError(t, err)
}

func TestDeleteByID(t *testing.T) {
	svc, m := getTestService(t)

	m.On("DeleteByID", 1).Return(nil)

	err := svc.DeleteByID(1)
	require.NoError(t, err)

	m.On("DeleteByID", 2).Return(sql.ErrNoRows)

	err = svc.DeleteByID(2)
	require.ErrorIs(t, err, utils.ErrNotFound)
}

func getTestService(t *testing.T) (company.Service, *mocker) {
	var repo company.Service
	m := &mocker{}

	go fxtest.New(
		fxtest.TB(t),
		fx.Options(
			configs.Module,
			logger.Module,
			gateways.Module,

			fx.Provide(
				func() companyRepo.Repository {
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

func (m *mocker) Create(c companyRepo.Company) (err error) {
	args := m.Called(c)
	return args.Error(0)
}

func (m *mocker) GetByID(id int) (c companyRepo.Company, err error) {
	args := m.Called(id)
	return args.Get(0).(companyRepo.Company), args.Error(1)
}

func (m *mocker) GetAll(f companyRepo.Filters) (companies []companyRepo.Company, err error) {
	args := m.Called(f)
	companies, _ = args.Get(0).([]companyRepo.Company)

	return companies, args.Error(1)
}

func (m *mocker) Update(c companyRepo.Company) (err error) {
	args := m.Called(c)
	return args.Error(0)
}

func (m *mocker) DeleteByID(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}
