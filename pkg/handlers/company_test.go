package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xm/configs"
	"xm/pkg/db"
	"xm/pkg/handlers"
	"xm/pkg/logger"
	"xm/pkg/repositories"
	"xm/pkg/repositories/company"
	companyService "xm/pkg/services/company"
	"xm/pkg/services/user"
	"xm/pkg/services/utils"

	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCreateCompany(t *testing.T) {
	h, m := getTestHandlerCompany(t)

	tests := []struct {
		name     string
		c        company.Company
		status   int
		expected string
	}{
		{
			name: "ok",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status:   200,
			expected: `{"code":200,"message":"OK","payload":"created"}`,
		},
		{
			name: "no name",
			c: company.Company{
				Code:    "code",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status:   400,
			expected: `{"code":400,"message":"Bad Request","payload":"no name provided"}`,
		},
		{
			name: "no code",
			c: company.Company{
				Name:    "name",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status:   400,
			expected: `{"code":400,"message":"Bad Request","payload":"no code provided"}`,
		},
		{
			name: "no country",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Website: "website",
				Phone:   "phone",
			},
			status:   400,
			expected: `{"code":400,"message":"Bad Request","payload":"no country provided"}`,
		},
		{
			name: "no website",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Phone:   "phone",
			},
			status:   400,
			expected: `{"code":400,"message":"Bad Request","payload":"no website provided"}`,
		},
		{
			name: "no phone",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Website: "website",
			},
			status:   400,
			expected: `{"code":400,"message":"Bad Request","payload":"no phone provided"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.On("Create", tt.c).Return(nil)

			cJson, _ := json.Marshal(tt.c)

			req := httptest.NewRequest("POST", "/company/create", bytes.NewBuffer(cJson))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.CreateCompany)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}

			if rr.Body.String() != tt.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expected)
			}
		})
	}
}

func TestGetCompanyByID(t *testing.T) {
	h, m := getTestHandlerCompany(t)

	tests := []struct {
		name   string
		c      company.Company
		status int
		err    error
	}{
		{
			name: "ok",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status: 200,
			err:    nil,
		},
		{
			name: "not found",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status: 404,
			err:    utils.ErrNotFound,
		},
		{
			name: "internal error",
			c: company.Company{
				Name:    "name",
				Code:    "code",
				Country: "country",
				Website: "website",
				Phone:   "phone",
			},
			status: 500,
			err:    errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.On("GetByID", 1).Return(tt.c, tt.err).Once()

			req := httptest.NewRequest("GET", "/company?id=1", nil)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetCompanyByID)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}
		})
	}
}

func TestGetAllCompanies(t *testing.T) {
	h, m := getTestHandlerCompany(t)

	tests := []struct {
		name    string
		c       []company.Company
		filters string
		status  int
		err     error
	}{
		{
			name:    "ok",
			c:       []company.Company{{}, {}},
			status:  200,
			filters: `{"limit":10}`,
			err:     nil,
		},
		{
			name:    "bad requst",
			c:       nil,
			status:  400,
			filters: ``,
			err:     nil,
		},
		{
			name:    "not found",
			c:       nil,
			status:  404,
			filters: `{"limit":12}`,
			err:     utils.ErrNotFound,
		},
		{
			name:    "internal error",
			c:       nil,
			status:  500,
			filters: `{"limit":13}`,
			err:     errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filters company.Filters
			json.NewDecoder(strings.NewReader(tt.filters)).Decode(&filters)

			m.On("GetAll", filters).Return(tt.c, tt.err).Once()

			req := httptest.NewRequest("POST", "/companies", strings.NewReader(tt.filters))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetAllCompanies)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}
		})
	}
}

func TestDeleteCompany(t *testing.T) {
	h, m := getTestHandlerCompany(t)

	tests := []struct {
		name    string
		queryID string
		id      int
		status  int
		err     error
	}{
		{
			id:      1,
			queryID: "?id=1",
			name:    "ok",
			status:  200,
			err:     nil,
		},
		{
			name:   "bad requst",
			status: 400,
			err:    nil,
		},
		{
			id:      2,
			queryID: "?id=2",
			name:    "not found",
			status:  404,
			err:     utils.ErrNotFound,
		},
		{
			id:      3,
			queryID: "?id=3",
			name:    "internal error",
			status:  500,
			err:     errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.On("DeleteByID", tt.id).Return(tt.err)

			req := httptest.NewRequest("DELETE", "/company"+tt.queryID, nil)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.DeleteCompany)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}
		})
	}
}

func TestUpdateCompany(t *testing.T) {
	h, m := getTestHandlerCompany(t)

	tests := []struct {
		name   string
		cmp    string
		status int
		err    error
	}{
		{
			name:   "ok",
			status: 200,
			cmp:    `{"name":"some name"}`,
			err:    nil,
		},
		{
			name:   "bad requst",
			cmp:    ``,
			status: 400,
			err:    nil,
		},
		{
			name:   "internal error",
			cmp:    `{"name": "error"}`,
			status: 500,
			err:    errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c company.Company
			json.NewDecoder(strings.NewReader(tt.cmp)).Decode(&c)

			m.On("Update", c).Return(tt.err).Once()

			req := httptest.NewRequest("PATCH", "/company", strings.NewReader(tt.cmp))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateCompany)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}
		})
	}
}

func getTestHandlerCompany(t *testing.T) (handlers.Handlers, *companyMocker) {
	var h handlers.Handlers
	m := &companyMocker{}

	go fxtest.New(
		fxtest.TB(t),
		fx.Options(
			configs.Module,
			logger.Module,
			handlers.Module,
			repositories.Module,
			db.Module,

			user.Module,
			fx.Provide(
				func() companyService.Service {
					return m
				},
			),
		),
		fx.Populate(&h),
	).Run()

	return h, m
}

type companyMocker struct {
	mock.Mock
}

func (m *companyMocker) Create(c company.Company) (err error) {
	args := m.Called(c)
	return args.Error(0)
}

func (m *companyMocker) GetByID(id int) (c company.Company, err error) {
	args := m.Called(id)
	return args.Get(0).(company.Company), args.Error(1)
}

func (m *companyMocker) GetAll(f company.Filters) (companies []company.Company, err error) {
	args := m.Called(f)
	companies, _ = args.Get(0).([]company.Company)

	return companies, args.Error(1)
}

func (m *companyMocker) Update(c company.Company) (err error) {
	args := m.Called(c)
	return args.Error(0)
}

func (m *companyMocker) DeleteByID(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}
