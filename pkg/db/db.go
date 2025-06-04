package db

import (
	"fmt"
	"time"

	"github.com/zgsm/go-webserver/config"
	"github.com/zgsm/go-webserver/i18n"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// InitDB 初始化数据库连接
func InitDB(cfg config.Database) error {
	var err error
	var dialector gorm.Dialector

	switch cfg.Type {
	// case "mysql":
	// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
	// 		cfg.User,
	// 		cfg.Password,
	// 		cfg.Host,
	// 		cfg.Port,
	// 		cfg.DBName,
	// 		cfg.TimeZone)
	// 	dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
			cfg.SSLMode,
			cfg.TimeZone)
		dialector = postgres.Open(dsn)
	// case "sqlite":
	// 	dialector = sqlite.Open(cfg.DBName)
	default:
		return fmt.Errorf("%s", i18n.Translate("db.unsupported_type", "", map[string]interface{}{"type": cfg.Type}))
	}

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	}

	// 连接数据库
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.Translate("db.connection.failed", "", nil), err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.Translate("db.connection.failed", "", nil), err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate 自动迁移模型
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("%s", i18n.Translate("db.not_initialized", "", nil))
	}

	return DB.AutoMigrate(models...)
}
