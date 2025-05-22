package main

import (
	"log"
	"os"

	"time"

	"github.com/zgsm/review-manager/i18n"
	"github.com/zgsm/review-manager/pkg/i18nlogger"

	"github.com/spf13/cobra"
	"github.com/zgsm/review-manager/config"
	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/pkg/db"
	"github.com/zgsm/review-manager/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: i18n.Translate("db.init.short", "", nil),
	Long:  i18n.Translate("db.init.long", "", nil),
	Run: func(cmd *cobra.Command, args []string) {
		// 加载配置
		if err := config.LoadConfigWithDefault(); err != nil {
			log.Fatalf(i18n.Translate("config.load.failed", "", nil)+": %v", err)
		}

		// 初始化日志
		if err := logger.InitLogger(config.GetConfig().Log); err != nil {
			log.Fatalf(i18n.Translate("logger.init.failed", "", nil)+": %v", err)
		}
		defer logger.Sync()

		// 初始化数据库连接
		if err := db.InitDB(*config.GetConfig()); err != nil {
			logger.Error(i18n.Translate("db.connection.init", "", nil), "error", err)
			i18nlogger.Error("db.init.failed", "", nil, "error", err)
			os.Exit(1)
		}

		// 执行数据库迁移
		logger.Info(i18n.Translate("db.migrate.start", "", nil))
		i18nlogger.Info("db.migrate.start", "", nil)

		// 注册所有需要迁移的模型
		logger.Info(i18n.Translate("db.model.register", "", nil))
		if err := db.AutoMigrate(
			&model.User{},
			&model.UserProfile{},
			// 在这里添加其他模型
		); err != nil {
			i18nlogger.Error("db.migrate.failed", "", nil, "error", err)
			os.Exit(1)
		}

		// 添加初始数据
		i18nlogger.Info("db.seed.start", "", nil)

		// 添加管理员用户
		logger.Info(i18n.Translate("db.admin.create", "", nil))
		if err := createAdminUser(); err != nil {
			i18nlogger.Error("user.admin.create.failed", "", nil, "error", err)
			os.Exit(1)
		}

		i18nlogger.Info("db.init.success", "", nil)
	},
}

// createAdminUser 创建管理员用户
func createAdminUser() error {
	// 检查管理员用户是否已存在
	logger.Info(i18n.Translate("db.check.admin_exists", "", nil))
	var count int64
	if err := db.GetDB().Model(&model.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	// 如果已存在管理员用户，则跳过
	if count > 0 {
		logger.Info(i18n.Translate("db.skip.admin_exists", "", nil))
		i18nlogger.Info("db.admin.exists", "", nil)
		return nil
	}

	// 生成密码哈希
	logger.Info(i18n.Translate("db.hash.generate", "", nil))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建管理员用户
	logger.Info(i18n.Translate("db.admin.create", "", nil))
	now := time.Now()
	admin := model.User{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  string(hashedPassword),
		Nickname:  i18n.Translate("user.admin.system", "", nil),
		Role:      "admin",
		Status:    1,
		LastLogin: &now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 保存到数据库
	logger.Info(i18n.Translate("db.save.db", "", nil))
	if err := db.GetDB().Create(&admin).Error; err != nil {
		return err
	}

	// 创建管理员资料
	logger.Info(i18n.Translate("db.profile.create", "", nil))
	adminProfile := model.UserProfile{
		UserID:    admin.ID,
		RealName:  i18n.Translate("user.admin.system", "", nil),
		Phone:     "13800138000",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 保存到数据库
	logger.Info(i18n.Translate("db.save.db", "", nil))
	if err := db.GetDB().Create(&adminProfile).Error; err != nil {
		return err
	}

	i18nlogger.Info("db.admin.created", "", nil, "username", admin.Username)
	return nil
}
