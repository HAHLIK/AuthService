package sqlite

import (
	"context"

	"github.com/HAHLIK/AuthService/sso/internal/domain/models"
)

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	return 1, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	return models.User{}, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	return true, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	return models.App{}, nil
}
