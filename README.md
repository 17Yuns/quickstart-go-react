# Go Viper 配置管理系统

基于 Viper 和 Gin 的配置管理系统，支持热更新和自动重启服务器。

## 功能特性

- ✅ **配置文件解析**: 使用 Viper 解析 YAML 配置文件
- ✅ **热更新**: 配置文件变化时自动重新加载
- ✅ **自动重启**: 配置变化后自动重启 Gin 服务器
- ✅ **线程安全**: 支持并发安全的配置访问
- ✅ **类型安全**: 使用结构体映射，避免字符串键值访问
- ✅ **优雅关闭**: 支持优雅关闭服务器
- ✅ **全局访问**: 在任何文件中都可以方便地访问配置

## 项目结构

```
.
├── config/
│   ├── types.go    # 配置结构体定义
│   └── config.go   # 配置管理器
├── server/
│   └── server.go   # 服务器管理
├── handlers/
│   └── example.go  # 示例处理器
├── config.yaml     # 配置文件
├── main.go         # 主程序
└── README.md       # 说明文档
```

## 配置文件格式

```yaml
system:
  name: "QuickStart-React"
  port: 3000
  host: "0.0.0.0"
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行程序

```bash
go run main.go
```

### 3. 测试热更新

修改 `config.yaml` 文件中的端口号或应用名称，保存后观察控制台输出，服务器会自动重启。

## API 接口

启动后可以访问以下接口：

- `GET /health` - 健康检查
- `GET /config` - 获取当前配置
- `GET /api/info` - 获取应用信息

## 在其他文件中使用配置

### 方法1: 直接获取配置对象

```go
package yourpackage

import "quickstart-react/config"

func someFunction() {
    cfg := config.GetConfig()
    
    // 使用配置
    appName := cfg.System.Name
    port := cfg.System.Port
    host := cfg.System.Host
}
```

### 方法2: 在 Gin 处理器中使用

```go
func yourHandler(c *gin.Context) {
    cfg := config.GetConfig()
    
    c.JSON(200, gin.H{
        "app_name": cfg.System.Name,
        "port":     cfg.System.Port,
    })
}
```

### 方法3: 使用便捷方法

```go
import "quickstart-go-react/config"

// 获取字符串值
appName := config.GetString("system.name")

// 获取整数值
port := config.GetInt("system.port")
```

## 扩展配置

### 1. 添加新的配置项

在 `config.yaml` 中添加新配置：

```yaml
system:
  name: "QuickStart-React"
  port: 3000
  host: "0.0.0.0"

database:
  host: "localhost"
  port: 5432
  name: "mydb"
  user: "admin"
  password: "password"
```

### 2. 更新配置结构体

在 `config/types.go` 中添加对应的结构体：

```go
type Config struct {
    System   SystemConfig   `yaml:"system"`
    Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Name     string `yaml:"name"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
}
```

### 3. 使用新配置

```go
cfg := config.GetConfig()
dbHost := cfg.Database.Host
dbPort := cfg.Database.Port
```

## 配置变化回调

你可以注册回调函数来监听配置变化：

```go
config.OnConfigChange(func(cfg *config.Config) {
    log.Printf("配置已更新: %+v", cfg)
    // 执行你的自定义逻辑
})
```

## 注意事项

1. **配置文件路径**: 程序会在当前目录和 `./config` 目录下查找 `config.yaml` 文件
2. **端口占用**: 修改端口后，旧端口会在服务器重启时自动释放
3. **并发安全**: `GetConfig()` 方法是线程安全的，可以在多个 goroutine 中安全使用
4. **配置副本**: `GetConfig()` 返回配置的副本，外部修改不会影响原始配置

## 常见问题

### Q: 如何添加环境变量支持？

A: 在 `config/config.go` 的 `Init()` 函数中添加：

```go
viper.AutomaticEnv()
viper.SetEnvPrefix("APP") // 环境变量前缀
```

### Q: 如何支持多种配置文件格式？

A: Viper 支持多种格式，只需修改 `SetConfigType()` 即可：

```go
viper.SetConfigType("json") // 或 "toml", "properties" 等
```

### Q: 如何禁用热更新？

A: 注释掉 `Init()` 函数中的 `viper.WatchConfig()` 相关代码即可。

## 性能说明

- 配置读取使用读写锁，读操作性能优异
- 配置变化时会创建新的配置副本，内存使用合理
- 服务器重启过程中会有短暂的服务中断（通常小于1秒）

## 许可证

MIT License