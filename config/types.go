package config

// Config 主配置结构体
type Config struct {
	System SystemConfig `yaml:"system" mapstructure:"system"`
	Log    LogConfig    `yaml:"log" mapstructure:"log"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Name string `yaml:"name" mapstructure:"name"` // 应用名称
	Port int    `yaml:"port" mapstructure:"port"` // 服务端口
	Host string `yaml:"host" mapstructure:"host"` // 服务主机地址
}

// LogConfig 日志配置
type LogConfig struct {
	Level         string `yaml:"level" mapstructure:"level"`                   // 日志级别
	Format        string `yaml:"format" mapstructure:"format"`                 // 日志格式
	OutputDir     string `yaml:"output_dir" mapstructure:"output_dir"`         // 日志输出目录
	MaxSize       int    `yaml:"max_size" mapstructure:"max_size"`             // 单个日志文件最大大小(MB)
	MaxBackups    int    `yaml:"max_backups" mapstructure:"max_backups"`       // 保留的旧日志文件数量
	MaxAge        int    `yaml:"max_age" mapstructure:"max_age"`               // 日志文件保留天数
	Compress      bool   `yaml:"compress" mapstructure:"compress"`             // 是否压缩旧日志文件
	ConsoleOutput bool   `yaml:"console_output" mapstructure:"console_output"` // 是否同时输出到控制台
}
