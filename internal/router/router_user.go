package router

import (
	"github.com/gin-gonic/gin"
	"monitor-apiserver/internal/controller/v1/user"
	"monitor-apiserver/internal/repository"
)

func setUserRouter(r *gin.RouterGroup, repo repository.Repository) {
	api := r.Group("/user", AuthToken())
	userController := user.NewUserController(repo)
	{
		api.GET("", userController.GetUserInfo())
	}
}
