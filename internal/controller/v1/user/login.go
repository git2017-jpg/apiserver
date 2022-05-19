package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"monitor-apiserver/internal/model"
	"monitor-apiserver/pkg/config"
	"monitor-apiserver/pkg/errors"
	"monitor-apiserver/pkg/errors/code"
	"monitor-apiserver/pkg/jwt"
	"monitor-apiserver/pkg/response"
	"monitor-apiserver/pkg/security"
	jtime "monitor-apiserver/pkg/time"
	"time"
)

func (u *UserController) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginReqParam := model.LoginReq{}
		if err := c.ShouldBind(&loginReqParam); err != nil {
			response.JSON(c, errors.WithCode(code.ValidateErr, err.Error()), nil)
			return
		}
		// 查询用户信息
		user, err := u.srv.Users().GetByMobile(context.TODO(), loginReqParam.Mobile)
		if err != nil {
			response.JSON(c, errors.Wrap(err, code.UserLoginErr, "登录失败，用户不存在"), nil)
			return
		}

		if !security.ValidatePassword(loginReqParam.Password, user.Password) {
			response.JSON(c, errors.WithCode(code.UserLoginErr, "登录失败，用户名、密码不匹配"), nil)
			return
		}
		// 生成jwt token
		expireAt := time.Now().Add(24 * 7 * time.Hour)
		claims := jwt.BuildClaims(expireAt, user.ID)
		token, err := jwt.GenToken(claims, config.GlobalConfig.JwtSecret)
		if err != nil {
			response.JSON(c, errors.Wrap(err, code.UserLoginErr, "生成用户授权token失败"), nil)
			return
		}
		response.JSON(c, nil, struct {
			Token    string         `json:"token"`
			ExpireAt jtime.JsonTime `json:"expire_at"`
		}{
			Token:    token,
			ExpireAt: jtime.JsonTime(expireAt),
		})
	}
}
