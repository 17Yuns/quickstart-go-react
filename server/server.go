package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"quickstart-go-react/config"
	"quickstart-go-react/logger"

	"github.com/gin-gonic/gin"
)

// Server Gin服务器封装
type Server struct {
	httpServer *http.Server
	engine     *gin.Engine
}

// New 创建新的服务器实例
func New() *Server {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Server{
		engine: engine,
	}
}

// SetupRoutes 设置路由
func (s *Server) SetupRoutes() {
	// 健康检查接口
	s.engine.GET("/health", func(c *gin.Context) {
		cfg := config.GetConfig()
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"app_name": cfg.System.Name,
			"port":     cfg.System.Port,
			"host":     cfg.System.Host,
			"time":     time.Now().Format("2006-01-02 15:04:05"),
		})
		logger.Infof("健康检查请求 - 客户端: %s", c.ClientIP())
	})

	// 获取配置信息接口
	s.engine.GET("/config", func(c *gin.Context) {
		cfg := config.GetConfig()
		c.JSON(http.StatusOK, gin.H{
			"config": cfg,
		})
		logger.Infof("配置信息请求 - 客户端: %s", c.ClientIP())
	})

	// 示例API接口
	s.engine.GET("/api/info", func(c *gin.Context) {
		cfg := config.GetConfig()
		c.JSON(http.StatusOK, gin.H{
			"message":  "这是一个示例API",
			"app_name": cfg.System.Name,
			"version":  "1.0.0",
		})
		logger.Infof("API信息请求 - 客户端: %s", c.ClientIP())
	})

	// 测试日志接口
	s.engine.GET("/test-log", func(c *gin.Context) {
		logger.Traceln("这是一条 trace 日志")
		logger.Debugln("这是一条 debug 日志")
		logger.Infoln("这是一条 info 日志")
		logger.Warnln("这是一条 warn 日志")
		logger.Errorln("这是一条 error 日志")

		c.JSON(http.StatusOK, gin.H{
			"message": "日志测试完成，请查看日志文件",
		})
	})

	// 静态文件服务（如果需要）
	s.engine.Static("/static", "./static")
}

// Start 启动服务器
func (s *Server) Start() error {
	cfg := config.GetConfig()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.System.Host, cfg.System.Port),
		Handler: s.engine,
	}

	logger.Infof("服务器启动中... 地址: http://%s:%d", cfg.System.Host, cfg.System.Port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	return nil
}

// Stop 优雅停止服务器
func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	logger.Infoln("正在优雅关闭服务器...")

	// 设置5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("服务器关闭失败: %v", err)
		return err
	}

	logger.Infoln("服务器已优雅关闭")
	return nil
}

// GetEngine 获取Gin引擎实例（用于添加自定义路由）
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}
