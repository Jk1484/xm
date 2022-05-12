package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xm/configs"
	"xm/gateways"
	"xm/pkg/db"
	"xm/pkg/handlers"
	"xm/pkg/logger"
	"xm/pkg/repositories"
	userRepo "xm/pkg/repositories/user"
	"xm/pkg/services/company"
	userService "xm/pkg/services/user"
	"xm/pkg/services/utils"

	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestSignUp(t *testing.T) {
	h, m := getTestHandler(t)

	tests := []struct {
		name     string
		data     string
		status   int
		expected string
		err      error
	}{
		{
			name:     "ok",
			data:     `{"username": "some", "password":"password"}`,
			expected: `{"code":200,"message":"OK","payload":"sign up completed"}`,
			status:   200,
			err:      nil,
		},
		{
			name:     "already exists",
			data:     `{"username": "some", "password":"password"}`,
			expected: `{"code":400,"message":"Bad Request","payload":"already registered"}`,
			status:   400,
			err:      utils.ErrAlreadyExists,
		},
		{
			name:     "bad request empty",
			data:     ``,
			expected: `{"code":400,"message":"Bad Request","payload":"bad credentials"}`,
			status:   400,
			err:      nil,
		},
		{
			name:     "no username",
			data:     `{"password": "somepass"}`,
			expected: `{"code":400,"message":"Bad Request","payload":"no username provided"}`,
			status:   400,
			err:      errors.New("some error"),
		},
		{
			name:     "bad password",
			data:     `{"username": "someuser"}`,
			expected: `{"code":400,"message":"Bad Request","payload":"password minimum length should be at least 8"}`,
			status:   400,
			err:      errors.New("some error"),
		},
		{
			name:     "internal err",
			data:     `{"username": "someuser", "password": "somepass"}`,
			expected: `{"code":500,"message":"Internal Server Error","payload":null}`,
			status:   500,
			err:      errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u userRepo.User
			err := json.NewDecoder(strings.NewReader(tt.data)).Decode(&u)

			if err == nil {
				m.On("Create", mock.Anything).Return(tt.err).Once()
			}

			req := httptest.NewRequest("POST", "/sign-up", strings.NewReader(tt.data))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.SignUp)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
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

func TestSignIn(t *testing.T) {
	h, m := getTestHandler(t)

	tests := []struct {
		name   string
		data   string
		user   userRepo.User
		status int
		err    error
	}{
		{
			name:   "ok",
			data:   `{"username": "some", "password":"password"}`,
			user:   userRepo.User{Username: "some", Password: "password"},
			status: 200,
			err:    nil,
		},
		{
			name:   "no username",
			data:   `{"username": "", "password":"password"}`,
			user:   userRepo.User{Username: "", Password: "password"},
			status: 400,
			err:    nil,
		},
		{
			name:   "no password",
			data:   `{"username": "some", "password":""}`,
			user:   userRepo.User{Username: "some", Password: "somepass"},
			status: 400,
			err:    nil,
		},
		{
			name:   "internal error",
			data:   `{"username": "some", "password":"password"}`,
			user:   userRepo.User{Username: "some", Password: "password"},
			status: 500,
			err:    errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u userRepo.User
			err := json.NewDecoder(strings.NewReader(tt.data)).Decode(&u)

			if err == nil {
				tt.user.Password, _ = handlers.HashPassword(tt.user.Password)
				m.On("GetByUsername", tt.user.Username).Return(&tt.user, tt.err).Once()
			}

			req := httptest.NewRequest("POST", "/sign-in", strings.NewReader(tt.data))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.SignIn)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.status {
				fmt.Println(status, tt.status, status == tt.status)
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.status)
			}
		})
	}
}

func getTestHandler(t *testing.T) (handlers.Handlers, *userMocker) {
	var h handlers.Handlers
	m := &userMocker{}

	go fxtest.New(
		fxtest.TB(t),
		fx.Options(
			configs.Module,
			logger.Module,
			handlers.Module,
			repositories.Module,
			db.Module,
			gateways.Module,

			company.Module,
			fx.Provide(
				func() userService.Service {
					return m
				},
			),
		),
		fx.Populate(&h),
	).Run()

	return h, m
}

type userMocker struct {
	mock.Mock
}

func (m *userMocker) Create(u *userRepo.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *userMocker) GetByUsername(username string) (*userRepo.User, error) {
	args := m.Called(username)
	v, _ := args.Get(0).(*userRepo.User)
	return v, args.Error(1)
}
