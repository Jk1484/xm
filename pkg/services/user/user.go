package user

import (
	"database/sql"
	"xm/pkg/repositories/user"
	"xm/pkg/services/utils"

	"github.com/lib/pq"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Service interface {
	Create(u *user.User) error
	GetByUsername(username string) (*user.User, error)
}

type service struct {
	userRepository user.Repository
}

type Params struct {
	fx.In
	UserRepository user.Repository
}

func New(p Params) Service {
	return &service{
		userRepository: p.UserRepository,
	}
}

func (s *service) Create(u *user.User) error {
	err := s.userRepository.Create(u)
	if err != nil {
		v, ok := err.(*pq.Error)
		if !ok {
			return err
		}

		if v.Code == "23505" {
			return utils.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (s *service) GetByUsername(username string) (*user.User, error) {
	u, err := s.userRepository.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	return u, nil
}
