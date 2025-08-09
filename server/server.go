package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"quickstart-go-react/config"
	"quickstart-go-react/logger"
	"quickstart-go-react/routes"

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
	routes.SetupRoutes(s.engine)
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
