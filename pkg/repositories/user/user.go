package user

import (
	"database/sql"
	"time"
	"xm/pkg/db"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Repository interface {
	Create(u *User) error
	GetByUsername(username string) (*User, error)
}

type repository struct {
	db *sql.DB
}

type Params struct {
	fx.In
	DB db.Database
}

func New(p Params) Repository {
	return &repository{
		db: p.DB.Connection(),
	}
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *repository) Create(u *User) error {

	query := `
		INSERT INTO users (
			username,
			password
		)
		VALUES($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, u.Username, u.Password).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetByUsername(username string) (*User, error) {
	var u User

	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
