// Created on 2021/5/4.
// @author tony
// email xmgtony@gmail.com
// description jwt中间件

package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"monitor-apiserver/pkg/config"
	"monitor-apiserver/pkg/errors"
	"monitor-apiserver/pkg/errors/code"
	"monitor-apiserver/pkg/jwt"
	"monitor-apiserver/pkg/response"
	"strings"
)

// 请求头的形式为 Authorization: Bearer token
const authorizationHeader = "Authorization"

//AuthToken 鉴权，验证用户token是否有效
func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getJwtFromHeader(c)
		if err != nil {
			response.JSON(c, errors.Wrap(err, code.RequireAuthErr, "invalid token"), nil)
			c.Abort()
			return
		}
		// 验证token是否正确
		claims, err := jwt.ParseToken(token, config.GlobalConfig.JwtSecret)
		if err != nil {
			response.JSON(c, errors.Wrap(err, code.RequireAuthErr, "invalid token"), nil)
			c.Abort()
			return
		}
		c.Set(config.UserID, claims.UserId)
		c.Next()
	}
}

func getJwtFromHeader(c *gin.Context) (string, error) {
	aHeader := c.Request.Header.Get(authorizationHeader)
	if len(aHeader) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	strs := strings.SplitN(aHeader, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", fmt.Errorf("token 不符合规则")
	}
	return strs[1], nil
}
