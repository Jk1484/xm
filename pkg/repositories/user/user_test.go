package user_test

import (
	"database/sql"
	"fmt"
	"testing"
	"xm/configs"
	"xm/pkg/db"
	"xm/pkg/logger"
	"xm/pkg/repositories"
	"xm/pkg/repositories/user"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreate(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	var u = user.User{
		Username: "Some",
		Password: "Pass",
	}
	err = repo.Create(&u)
	require.NoError(t, err)

	err = repo.Create(&u)
	require.Error(t, err)

	u2, err := repo.GetByUsername(u.Username)
	require.NoError(t, err)
	require.Equal(t, u.Password, u2.Password)
}

func TestGetByUsername(t *testing.T) {
	repo, err := getTestRepo(t)
	require.NoError(t, err)

	var u = user.User{
		Username: "Some",
		Password: "Pass",
	}
	err = repo.Create(&u)
	require.NoError(t, err)

	u = user.User{
		Username: "some2",
		Password: "pass2",
	}
	err = repo.Create(&u)
	require.NoError(t, err)

	u2, err := repo.GetByUsername(u.Username)
	require.NoError(t, err)
	require.Equal(t, u.Password, u2.Password)
}

func getTestRepo(t *testing.T) (user.Repository, error) {
	var repo user.Repository
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

	err := createUsersTable(dbConn)
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

type mock struct {
	db *sql.DB
}

func (m *mock) Connection() *sql.DB {
	return m.db
}

func (m *mock) CloseConnection() error {
	return m.db.Close()
}

func createUsersTable(db *sql.DB) (err error) {
	query := `
		DROP TABLE IF EXISTS users
	`

	db.Exec(query)

	query = `
		CREATE TABLE users (
			id serial primary key,
			username varchar(100) unique,
			password varchar(100),
			created_at timestamp default now(),
			updated_at timestamp default now()
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		return
	}

	return
}
