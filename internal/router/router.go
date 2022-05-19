package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"monitor-apiserver/internal/repository"
	"monitor-apiserver/pkg/config"
	"monitor-apiserver/pkg/database"
	"monitor-apiserver/pkg/log"
	"monitor-apiserver/pkg/middleware"
	"monitor-apiserver/pkg/shutdown"
	"monitor-apiserver/pkg/utils"
	"net/http"
	"time"
)

type Application struct {
	config *config.Config
}

func NewApplication(config *config.Config) *Application {
	return &Application{config: config}
}

func (a *Application) Run() error {
	// 创建数据库链接，使用默认的实现方式
	ds := database.NewDefaultMysql(a.config.DBConfig)
	repo := repository.NewRepository(ds)
	// 设置gin启动模式，必须在创建gin实例之前
	gin.SetMode(a.config.Mode)
	g := a.genRouter(repo)

	server := http.Server{
		Addr:    a.config.Port,
		Handler: g,
	}

	// health check
	go func() {
		if err := a.ping(); err != nil {
			log.Error("server no response")
		}
		log.Infof("server started success! port: %s", a.config.Port)
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("http server startup err: ", err.Error())
		}
		log.Infof("server start failed on port %s", a.config.Port)
	}()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Errorf("server shutdown err：", err.Error())
			}
		},

		// 关闭 db
		func() {
			if ds != nil {
				ds.Close()
			}
		},

		// 关闭 cache
		func() {},

		// 关闭 cron Server
		func() {},
	)
	return nil
}

func (a *Application) genRouter(repo repository.Repository) *gin.Engine {
	g := gin.New()
	// 使用中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	//g.Use(middleware.Logging())
	//g.Use(middleware.AccessLog())
	g.Use(middleware.RequestID())
	g.Use(middleware.Metrics(a.config.AppName))
	g.Use(middleware.Timeout(3 * time.Second))

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "404 not found!")
	})

	// 405 Handler
	g.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "405 not method!")
	})

	// HealthCheck 健康检查路由
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "Up",
			"app":      a.config.AppName,
			"hostname": utils.GetHostname(),
		})
	})

	// metrics router 可以在 prometheus 中进行监控
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// apiV1 业务路由
	apiV1 := g.Group("/v1")
	apiV1.Use()
	{
		setUserRouter(apiV1, repo)
	}
	return g
}

func (a *Application) ping() error {
	seconds := 1
	url := "http://127.0.0.1" + a.config.Port + "/health"
	for i := 0; i < a.config.MaxPingCount; i++ {
		resp, err := http.Get(url)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		log.Infof("等待服务在线, 已等待 %d 秒，最多等待 %d 秒", seconds, a.config.MaxPingCount)
		time.Sleep(time.Second * 1)
		seconds++
	}
	return fmt.Errorf("服务启动失败，端口 %s", a.config.Port)
}
