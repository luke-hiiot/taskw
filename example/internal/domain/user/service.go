package user

import (
	"context"

	"{{.Module}}/ent"
	"{{.Module}}/ent/user"
	. "{{.Module}}/internal/pkg/errors"
	"{{.Module}}/internal/pkg/logger"
)

type Service interface {
	GetUser(ctx context.Context, id int) (*ent.User, error)
	Repository
}

type ServiceImpl struct {
	Repository
	logger *logger.Logger
}

func ProvideService(repo Repository, logger *logger.Logger) Service {
	return &ServiceImpl{Repository: repo, logger: logger}
}

func (s *ServiceImpl) GetUser(ctx context.Context, id int) (*ent.User, error) {
	u, err := s.FindOne(ctx, user.IDEQ(id))
	if err != nil {
		return nil, ErrNotFound
	}
	return u, nil
}
