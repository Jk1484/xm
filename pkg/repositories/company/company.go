package company

import (
	"database/sql"
	"strconv"
	"time"
	"xm/pkg/db"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Repository interface {
	Create(c Company) (err error)
	GetByID(id int) (c Company, err error)
	GetAll(f Filters) (companies []Company, err error)
	Update(c Company) (err error)
	DeleteByID(id int) (err error)
}

type repository struct {
	db *sql.DB
}

type Params struct {
	fx.In
	DB db.Database
}

type Company struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Country   string    `json:"country"`
	Website   string    `json:"website"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func New(p Params) Repository {
	return &repository{
		db: p.DB.Connection(),
	}
}

func (r *repository) Create(c Company) (err error) {
	query := `
		INSERT INTO companies(name, code, country, website, phone)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(query, c.Name, c.Code, c.Country, c.Website, c.Phone)
	if err != nil {
		return
	}

	return
}

func (r *repository) GetByID(id int) (c Company, err error) {
	query := `
		SELECT
			id, name, code, country, website, phone, status, created_at, updated_at
		FROM companies
		WHERE id = $1 AND status = 'active'
	`

	err = r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Code, &c.Country, &c.Website, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return
	}

	return
}

type Filters struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

func (r *repository) GetAll(f Filters) (companies []Company, err error) {
	query := `
		SELECT
			id, name, code, country, website, phone, status, created_at, updated_at
		FROM companies
		WHERE status = 'active'
	`

	cnt := 1
	var values []interface{}

	if f.ID != 0 {
		query += ` AND id = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.ID)
	}

	if f.Name != "" {
		query += ` AND name = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.Name)
	}

	if f.Code != "" {
		query += ` AND code = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.Code)
	}

	if f.Country != "" {
		query += ` AND country = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.Country)
	}

	if f.Website != "" {
		query += ` AND website = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.Website)
	}

	if f.Phone != "" {
		query += ` AND phone = $` + strconv.Itoa(cnt)
		cnt++

		values = append(values, f.Phone)
	}

	query += ` LIMIT $` + strconv.Itoa(cnt)
	cnt++
	values = append(values, f.Limit)

	query += ` OFFSET $` + strconv.Itoa(cnt)
	cnt++
	values = append(values, f.Offset)

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return
	}

	for rows.Next() {
		var c Company
		err = rows.Scan(&c.ID, &c.Name, &c.Code, &c.Country, &c.Website, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}

		companies = append(companies, c)
	}

	if len(companies) == 0 {
		return nil, sql.ErrNoRows
	}

	return
}

func (r *repository) Update(c Company) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	query := `
		UPDATE companies
		SET
			name = COALESCE(NULLIF($1, ''), name), code = COALESCE(NULLIF($2, ''), code),
			country = COALESCE(NULLIF($3, ''), country), website = COALESCE(NULLIF($4, ''), website),
			phone = COALESCE(NULLIF($5, ''), phone)
		WHERE id = $6
	`

	_, err = tx.Exec(query, c.Name, c.Code, c.Country, c.Website, c.Phone, c.ID)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	return tx.Commit()
}

func (r *repository) DeleteByID(id int) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	query := `
		UPDATE companies
		SET status = 'deleted'
		WHERE id = $1 AND status != 'deleted'
	`

	res, err := tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return
	}

	if cnt == 0 {
		_ = tx.Rollback()
		return sql.ErrNoRows
	}

	return tx.Commit()
}
