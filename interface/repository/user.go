package repository

import (
	"context"

	"github.com/trewanek/repository-with-tx/model"
)

type IUserRepository interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	Find(ctx context.Context, userID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userID string) error
}
