package company

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"xm/gateways/nats"
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
	natsGateway       nats.Gateway
}

type Params struct {
	fx.In
	CompanyRepository company.Repository
	NATSGateway       nats.Gateway
}

func New(p Params) Service {
	return &service{
		companyRepository: p.CompanyRepository,
		natsGateway:       p.NATSGateway,
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
	err = s.companyRepository.Update(c)
	if err != nil {
		return
	}

	m, _ := json.Marshal(c)
	s.natsGateway.GetConnection().Publish("company_update", m)

	return
}

func (s *service) DeleteByID(id int) (err error) {
	err = s.companyRepository.DeleteByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrNotFound
		}

		return
	}

	s.natsGateway.GetConnection().Publish("company_delete", []byte(strconv.Itoa(id)))

	return
}
