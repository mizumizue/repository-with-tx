package rdb

import (
	"context"

	"github.com/trewanek/repository-with-tx/model"
)

type UserMySqlRepository struct {
	handler SqlHandler
}

func NewUserMySqlRepository(handler SqlHandler) *UserMySqlRepository {
	return &UserMySqlRepository{handler: handler}
}

func (repo *UserMySqlRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	rows, err := repo.handler.query(ctx, "SELECT * FROM users;")
	if err != nil {
		return nil, err
	}

	res, err := repo.handler.rowsScan(ctx, rows, model.User{})
	if err != nil {
		return nil, err
	}

	return res.([]*model.User), nil
}

func (repo *UserMySqlRepository) Find(ctx context.Context, userID string) (*model.User, error) {
	dst := new(model.User)
	err := repo.handler.get(ctx, dst, "SELECT * FROM users WHERE user_id = ?;", userID)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (repo *UserMySqlRepository) Create(ctx context.Context, user *model.User) error {
	_, err := repo.handler.exec(
		ctx,
		"INSERT INTO users(user_name, email, telephone) VALUES(?, ?, ?);",
		user.UserName,
		user.Email,
		user.Telephone,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserMySqlRepository) Update(ctx context.Context, user *model.User) error {
	_, err := repo.handler.exec(
		ctx,
		"UPDATE users SET user_name = ?, email = ?, telephone = ? WHERE user_id = ?;",
		user.UserName,
		user.Email,
		user.Telephone,
		user.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserMySqlRepository) Delete(ctx context.Context, userID string) error {
	_, err := repo.handler.exec(
		ctx,
		"DELETE FROM users WHERE user_id = ?;",
		userID,
	)
	if err != nil {
		return err
	}
	return nil
}
