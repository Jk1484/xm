package company

import (
	"xm/pkg/repositories/company"

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
	return
}

func (s *service) GetByID(id int) (c company.Company, err error) {
	return
}

func (s *service) GetAll(f company.Filters) (companies []company.Company, err error) {
	return
}

func (s *service) Update(c company.Company) (err error) {
	return
}

func (s *service) DeleteByID(id int) (err error) {
	return
}
