package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config      *Config
	mutex       sync.RWMutex
	initialized bool
	initMutex   sync.Mutex
	callbacks   []func(*Config)
)

// Init 初始化配置管理器
func Init() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置到结构体
	if err := loadConfig(); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	// 启用配置文件监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("配置文件发生变化: %s", e.Name)

		// 添加短暂延迟，避免文件写入过程中的读取
		time.Sleep(100 * time.Millisecond)

		if err := loadConfig(); err != nil {
			logrus.Errorf("重新加载配置失败: %v", err)
			return
		}

		mutex.RLock()
		currentConfig := config
		mutex.RUnlock()

		logrus.Infof("配置重新加载成功，新配置: %+v", currentConfig)

		for _, callback := range callbacks {
			callback(currentConfig)
		}
	})

	initialized = true
	return nil
}

func ensureInit() {
	if !initialized {
		if err := Init(); err != nil {
			logrus.Fatalf("自动初始化配置失败: %v", err)
		}
	}
}

func loadConfig() error {
	var newConfig Config
	if err := viper.Unmarshal(&newConfig); err != nil {
		return err
	}

	mutex.Lock()
	config = &newConfig
	mutex.Unlock()

	return nil
}

func GetConfig() *Config {
	ensureInit() // 自动初始化

	mutex.RLock()
	defer mutex.RUnlock()

	configCopy := *config
	return &configCopy
}

func OnConfigChange(callback func(*Config)) {
	callbacks = append(callbacks, callback)
}
