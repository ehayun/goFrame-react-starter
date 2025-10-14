package repository

import (
	"context"
	"time"

	"tzlev/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.InsertedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := g.DB().Model("users").Ctx(ctx).Insert(user)
	return err
}

func (r *UserRepository) FindByZehut(ctx context.Context, zehut string) (*model.User, error) {
	var user model.User
	err := g.DB().Model("users").Ctx(ctx).
		Where("zehut = ?", zehut).
		Scan(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := g.DB().Model("users").Ctx(ctx).
		Where("email = ?", email).
		Scan(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	_, err := g.DB().Model("users").Ctx(ctx).
		Where("zehut = ?", user.Zehut).
		Update(user)
	return err
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]model.User, error) {
	var users []model.User
	err := g.DB().Model("users").Ctx(ctx).
		Offset(offset).
		Limit(limit).
		Scan(&users)

	return users, err
}
