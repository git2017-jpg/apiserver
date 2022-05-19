package service

import (
	"context"
	"monitor-apiserver/internal/model"
	"monitor-apiserver/internal/repository"
	"monitor-apiserver/pkg/errors"
	"monitor-apiserver/pkg/errors/code"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	GetByName(ctx context.Context, name string) (*model.User, error)
	GetById(ctx context.Context, uid int64) (*model.User, error)
	GetByMobile(ctx context.Context, mobile string) (*model.User, error)
}

type userService struct {
	repo repository.Repository
}

func newUserService(svc *service) UserService {
	return &userService{repo: svc.repo}
}

// GetByName 通过用户名 查找用户
func (us *userService) GetByName(ctx context.Context, name string) (*model.User, error) {
	if len(name) == 0 {
		return nil, errors.WithCode(code.ValidateErr, "用户名称不能为空")
	}
	return us.repo.Users().GetUserByName(ctx, name)
}

// GetById 根据用户ID查找用户
func (us *userService) GetById(ctx context.Context, uid int64) (*model.User, error) {
	return us.repo.Users().GetUserById(ctx, uid)
}

// GetByMobile 根据用户手机号查询
func (us *userService) GetByMobile(ctx context.Context, mobile string) (*model.User, error) {
	// 认为handler层对service层入参都是合法的，除了业务上的校验，service层不校验入参合规性
	return us.repo.Users().GetUserByMobile(ctx, mobile)
}
