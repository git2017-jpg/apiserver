package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"monitor-apiserver/internal/repository"
	"monitor-apiserver/internal/service"
	"monitor-apiserver/pkg/config"
	"monitor-apiserver/pkg/errors"
	"monitor-apiserver/pkg/errors/code"
	"monitor-apiserver/pkg/response"
)

type UserController struct {
	srv service.Service
}

func NewUserController(repo repository.Repository) *UserController {
	return &UserController{srv: service.NewService(repo)}
}

func (uh *UserController) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt64(config.UserID)
		user, err := uh.srv.Users().GetById(context.TODO(), uid)
		if err != nil {
			response.JSON(c, errors.Wrap(err, code.NotFoundErr, "用户信息为空"), nil)
		} else {
			response.JSON(c, nil, user)
		}
	}
}
