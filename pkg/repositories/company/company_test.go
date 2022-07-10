package company_test

import (
	"database/sql"
	"fmt"
	"testing"
	"xm/configs"
	"xm/pkg/db"
	"xm/pkg/logger"
	"xm/pkg/repositories"
	"xm/pkg/repositories/company"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreate(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	comp := company.Company{
		Name: "name",
		Code: "code",
	}

	err = repo.Create(comp)
	require.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	comp := company.Company{
		Name: "name",
		Code: "code",
	}

	err = repo.Create(comp)
	require.NoError(t, err)

	c, err := repo.GetByID(1)
	require.NoError(t, err)

	require.Equal(t, comp.Code, c.Code)
}

func TestGetAll(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	comp := company.Company{
		Name: "name",
		Code: "code",
	}

	err = repo.Create(comp)
	require.NoError(t, err)

	err = repo.Create(comp)
	require.NoError(t, err)

	err = repo.Create(company.Company{
		Name: "other",
		Code: "code",
	})
	require.NoError(t, err)

	c, err := repo.GetAll(company.Filters{Limit: 10})
	require.NoError(t, err)

	require.Equal(t, 3, len(c))

	c, err = repo.GetAll(company.Filters{Name: "other", Limit: 10})
	require.NoError(t, err)

	require.Equal(t, 1, len(c))

	c, err = repo.GetAll(company.Filters{Name: "not existing", Limit: 10})
	require.Equal(t, sql.ErrNoRows, err)

	require.Equal(t, 0, len(c))
}

func TestUpdate(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	c := company.Company{
		Name: "name1",
		Code: "ABC",
	}
	err = repo.Create(c)
	require.NoError(t, err)

	err = repo.Update(company.Company{
		ID:   1,
		Code: "EFG",
	})
	require.NoError(t, err)

	c2, err := repo.GetByID(1)
	require.NoError(t, err)

	require.Equal(t, c2.Code, "EFG")
}

func TestDeleteByID(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	c := company.Company{
		Name: "name1",
		Code: "ABC",
	}
	err = repo.Create(c)
	require.NoError(t, err)

	err = repo.DeleteByID(1)
	require.NoError(t, err)

	_, err = repo.GetByID(1)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

type mock struct {
	db *sql.DB
}

func getTestRepo(t *testing.T) (company.Repository, error) {
	var repo company.Repository
	var dbConn *sql.DB

	go fxtest.New(
		fxtest.TB(t),
		fx.Options(
			configs.Module,
			logger.Module,

			fx.Provide(
				func(
					p struct {
						fx.In
						Configs configs.Configs
					},
				) db.Database {
					dbConn = prepare(p.Configs)

					return &mock{
						db: dbConn,
					}
				},
			),

			repositories.Module,
		),
		fx.Populate(&repo),
	).Run()

	err := createCompaniesTable(dbConn)

	return repo, err
}

func prepare(cfg configs.Configs) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=testdb sslmode=disable",
		"testdb", "5432", cfg.Peek().Database.User, cfg.Peek().Database.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func (m *mock) Connection() *sql.DB {
	return m.db
}

func (m *mock) CloseConnection() error {
	return m.db.Close()
}

func createCompaniesTable(db *sql.DB) (err error) {
	query := `
		DROP TABLE IF EXISTS companies;
	`

	db.Exec(query)

	query = `
		CREATE TABLE companies(
			id serial primary key,
			name varchar(100) not null,
			code varchar(50) not null,
			country varchar(20) not null,
			website varchar(100) not null,
			phone varchar(50) not null,
			status varchar(20) not null default 'active',
			created_at timestamp not null default now(),
			updated_at timestamp not null default now()
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		return
	}

	return
}
