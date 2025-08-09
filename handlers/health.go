package handlers

import (
	"net/http"
	"time"

	"quickstart-go-react/config"
	"quickstart-go-react/logger"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器实例
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查接口
func (h *HealthHandler) Check(c *gin.Context) {
	cfg := config.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"app_name": cfg.System.Name,
		"port":     cfg.System.Port,
		"host":     cfg.System.Host,
		"time":     time.Now().Format("2006-01-02 15:04:05"),
	})
	logger.Infof("健康检查请求 - 客户端: %s", c.ClientIP())
}

// GetConfig 获取配置信息接口
func (h *HealthHandler) GetConfig(c *gin.Context) {
	cfg := config.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"config": cfg,
	})
	logger.Infof("配置信息请求 - 客户端: %s", c.ClientIP())
}
