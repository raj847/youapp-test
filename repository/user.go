package repository

import (
	"context"
	"youapp/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Register(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var res entity.User
	err := r.db.WithContext(ctx).Table("users").Where("username = ?", username).Find(&res).Error
	if err != nil {
		return entity.User{}, err
	}

	return res, nil
}
