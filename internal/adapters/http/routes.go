package http

import (
	"segmentation-service/pkg/infra/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRouter(a *Adapter, r *gin.Engine) error {
	log := logger.Get()
	log.Info("initializing handlers and routes...")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.Use(sloggin.New(log))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g := r.Group("/api/v1")
	{
		g.POST("/createSegment", a.createSegment)
		g.DELETE("/deleteSegment", a.deleteSegment)
		g.POST("/updateUserSegments/:userID", a.updateSegments)
		g.GET("/getUserSegments/:userID", a.getSegments)
		g.GET("/getReport/:period", a.getReport)
		g.GET("/getUserReport/:period/:userID", a.getUserReport)
	}
	return nil
}
