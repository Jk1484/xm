package company

import (
	"database/sql"
	"xm/pkg/repositories/company"
	"xm/pkg/services/utils"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Service interface {
	Create(c company.Company) (err error)
	GetByID(id int) (c company.Company, err error)
	GetAll(f company.Filters) (companies []company.Company, err error)
	Update(c company.Company) (err error)
	DeleteByID(id int) (err error)
}

type service struct {
	companyRepository company.Repository
}

type Params struct {
	fx.In
	CompanyRepository company.Repository
}

func New(p Params) Service {
	return &service{
		companyRepository: p.CompanyRepository,
	}
}

func (s *service) Create(c company.Company) (err error) {
	return s.companyRepository.Create(c)
}

func (s *service) GetByID(id int) (c company.Company, err error) {
	c, err = s.companyRepository.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return company.Company{}, utils.ErrNotFound
		}

		return
	}

	return
}

func (s *service) GetAll(f company.Filters) (companies []company.Company, err error) {
	companies, err = s.companyRepository.GetAll(f)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}

		return
	}

	return
}

func (s *service) Update(c company.Company) (err error) {
	return s.companyRepository.Update(c)
}

func (s *service) DeleteByID(id int) (err error) {
	err = s.companyRepository.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrNotFound
		}

		return
	}

	return
}
