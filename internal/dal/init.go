package dal

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"yikou-ai-go-teach/config"
	"yikou-ai-go-teach/internal/dal/query"
)

func InitDB(config *config.Config) *gorm.DB {
	if config == nil {
		panic(fmt.Errorf("配置加载失败"))
	}

	dsn := config.Database.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败: %w", err))
	}
	query.SetDefault(db)

	return db
}
