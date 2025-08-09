package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"quickstart-go-react/config"
	"quickstart-go-react/logger"
	"quickstart-go-react/server"
)

// App 应用程序结构体
type App struct {
	server       *server.Server
	serverMux    sync.Mutex
	restartCh    chan struct{}
	shutdownCh   chan struct{}
	lastRestart  time.Time
	restartMutex sync.Mutex
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{
		restartCh:  make(chan struct{}, 1),
		shutdownCh: make(chan struct{}, 1),
	}
}

// Start 启动应用
func (app *App) Start() error {
	// 初始化配置管理器
	if err := config.Init(); err != nil {
		return err
	}

	// 初始化日志系统
	if err := logger.Init(); err != nil {
		logger.Errorf("日志系统初始化失败: %v", err)
		return err
	}

	// 注册配置变化回调
	config.OnConfigChange(func(cfg *config.Config) {
		app.restartMutex.Lock()
		defer app.restartMutex.Unlock()

		// 防抖机制：如果距离上次重启不到1秒，则跳过
		if time.Since(app.lastRestart) < time.Second {
			logger.Warnln("重启过于频繁，跳过本次重启")
			return
		}

		logger.Infoln("检测到配置变化，准备重启服务器...")
		app.lastRestart = time.Now()

		select {
		case app.restartCh <- struct{}{}:
		default:
			logger.Warnln("重启信号已在队列中，跳过本次重启")
		}
	})

	// 启动服务器
	app.startServer()

	// 监听系统信号
	go app.handleSignals()

	// 监听重启信号
	go app.handleRestart()

	// 阻塞主线程
	<-app.shutdownCh
	logger.Infoln("应用程序退出")
	return nil
}

// startServer 启动服务器
func (app *App) startServer() {
	app.serverMux.Lock()
	defer app.serverMux.Unlock()

	// 创建新的服务器实例
	app.server = server.New()
	app.server.SetupRoutes()

	// 在新的goroutine中启动服务器
	go func() {
		if err := app.server.Start(); err != nil {
			logger.Errorf("服务器运行错误: %v", err)
		}
	}()

	// 等待一小段时间确保服务器启动
	time.Sleep(100 * time.Millisecond)
	cfg := config.GetConfig()
	logger.Infof("服务器启动成功 - %s 运行在 http://%s:%d",
		cfg.System.Name, cfg.System.Host, cfg.System.Port)
}

// stopServer 停止服务器
func (app *App) stopServer() {
	app.serverMux.Lock()
	defer app.serverMux.Unlock()

	if app.server != nil {
		if err := app.server.Stop(); err != nil {
			logger.Errorf("停止服务器时出错: %v", err)
		}
		app.server = nil
	}
}

// restartServer 重启服务器
func (app *App) restartServer() {
	logger.Infoln("开始重启服务器...")

	// 停止当前服务器
	app.stopServer()

	// 等待一小段时间确保端口释放
	time.Sleep(500 * time.Millisecond)

	// 启动新服务器
	app.startServer()

	logger.Infoln("服务器重启完成")
}

// handleRestart 处理重启信号
func (app *App) handleRestart() {
	for {
		select {
		case <-app.restartCh:
			app.restartServer()
		case <-app.shutdownCh:
			return
		}
	}
}

// handleSignals 处理系统信号
func (app *App) handleSignals() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	logger.Infof("收到信号: %v，开始优雅关闭...", sig)

	// 停止服务器
	app.stopServer()

	// 发送关闭信号
	close(app.shutdownCh)
}

func main() {
	logger.Infoln("启动应用程序...")

	app := NewApp()
	if err := app.Start(); err != nil {
		logger.Fatalf("应用启动失败: %v", err)
	}
}
