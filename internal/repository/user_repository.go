package repository

import (
	"context"
	"monitor-apiserver/internal/model"
	"monitor-apiserver/pkg/database"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserById(ctx context.Context, uid int64) (*model.User, error)
	GetUserByMobile(ctx context.Context, mobile string) (*model.User, error)
}

type userRepository struct {
	db database.IDataSource
}

func newUserRepository(repo repository) UserRepository {
	return &userRepository{db: repo.db}
}

func (ur *userRepository) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Master().Where("name = ?", name).Find(user).Error
	return user, err
}

func (ur *userRepository) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Master().Where("id = ?", uid).Find(user).Error
	return user, err
}

func (ur *userRepository) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Master().
		Where("mobile = ?", mobile).
		Where("enabled_status = 1").
		First(user).Error
	return user, err
}
