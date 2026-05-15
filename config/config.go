package config

import (
	"flag"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	AI       AIConfig       `mapstructure:"ai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int    `mapstructure:"port"`         // 服务端口
	ContextPath string `mapstructure:"context_path"` // api路径前缀
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // 数据库地址
	Port     int    `mapstructure:"port"`     // 数据库端口
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	Database string `mapstructure:"database"` // 数据库名
}

// AIConfig AI服务配置
type AIConfig struct {
	APIKey  string `mapstructure:"api_key"`  // API密钥
	Model   string `mapstructure:"model"`    // 模型名称
	BaseURL string `mapstructure:"base_url"` // API基础URL
}

// GetProjectRootPath 获取项目根路径
// 通过 runtime.Caller 获取当前文件的路径，然后向上查找 go.mod 文件所在目录
func GetProjectRootPath() (string, error) {
	// 获取当前文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}

	// 从当前文件路径向上查找，直到找到 go.mod 文件
	dir := filepath.Dir(filename)
	for {
		// 检查当前目录是否存在 go.mod 文件
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		// 向上一级目录
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// 已经到达根目录，仍未找到 go.mod
			return "", fmt.Errorf("未找到项目根路径（找不到 go.mod 文件）")
		}
		dir = parentDir
	}
}

// InitConfig 初始化配置
// env 参数用于指定配置文件后缀，如 "local" 会读取 config-local.yaml
func InitConfig() *Config {
	// 解析命令行参数
	env := flag.String("env", "", "运行环境，如 local, dev, test, prod")
	flag.Parse()

	// 获取项目根路径
	rootPath, err := GetProjectRootPath()
	if err != nil {
		panic(fmt.Errorf("获取项目根路径失败: %w", err))
	}

	// 拼接配置文件目录路径
	configPath := filepath.Join(rootPath, "config")

	// 确定配置文件名称
	configName := "config"
	if *env != "" {
		configName = fmt.Sprintf("config-%s", *env)
	}

	// 设置配置文件名和路径
	viper.SetConfigName(configName) // 配置文件名称
	viper.SetConfigType("yml")      // 配置文件类型
	viper.AddConfigPath(configPath) // 配置文件路径

	// 读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	logger.Infof("配置文件路径: %s\n", viper.ConfigFileUsed())

	// 解析配置到结构体
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("解析配置失败: %w", err))
	}
	return cfg
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
