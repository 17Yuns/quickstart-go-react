package routes

import (
	"quickstart-go-react/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(engine *gin.Engine) {
	// 创建处理器实例
	healthHandler := handlers.NewHealthHandler()

	// 健康检查路由
	engine.GET("/health", healthHandler.Check)
	engine.GET("/config", healthHandler.GetConfig)

}
