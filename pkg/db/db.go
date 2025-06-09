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

// InitDB initialize database connection
func InitDB(cfg config.Database) error {
	if !cfg.Enabled {
		return nil
	}

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

	// Configure GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
		DisableForeignKeyConstraintWhenMigrating: true, // Disable foreign key constraints
	}

	// Connect to database
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		_ = CloseDB()
		return fmt.Errorf("%s: %w", i18n.Translate("db.connection.failed", "", nil), err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		return fmt.Errorf("%s: %w", i18n.Translate("db.connection.failed", "", nil), err)
	}

	// Set max idle connections
	sqlDB.SetMaxIdleConns(10)
	// Set max open connections
	sqlDB.SetMaxOpenConns(100)
	// Set connection max lifetime
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

// GetDB get database connection
func GetDB() *gorm.DB {
	if DB == nil {
		panic("database is not initialized or disabled")
	}
	return DB
}

// CloseDB close database connection
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

// AutoMigrate auto migrate models
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("%s", i18n.Translate("db.not_initialized_or_disabled", "", nil))
	}

	return DB.AutoMigrate(models...)
}
