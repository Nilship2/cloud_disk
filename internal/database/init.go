// internal/database/init.go
package database

import (
	"fmt"
	"time"

	"cloud-disk/internal/config"
	"cloud-disk/internal/model/entity"
	"cloud-disk/pkg/logger"

	"github.com/glebarez/sqlite" // 纯 Go 实现的 SQLite 驱动
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 根据驱动类型选择数据库
	switch cfg.Driver {
	case "sqlite":
		// 使用纯 Go SQLite（不需要 CGO）
		db, err = gorm.Open(sqlite.Open("./cloud_disk.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			PrepareStmt: true,
		})
		if err == nil {
			logger.Info("SQLite database connected successfully")
		}

	case "mysql":
		// 使用 MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			cfg.Charset,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			PrepareStmt: true,
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
		if err == nil {
			logger.Info("MySQL database connected successfully",
				zap.String("host", cfg.Host),
				zap.Int("port", cfg.Port),
				zap.String("database", cfg.Database))
		}

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	/*获取通用数据库对象（仅 MySQL 需要）
	if cfg.Driver == "mysql" {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get sql.DB: %w", err)
		}
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}*/

	DB = db
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	logger.Info("Starting database auto migration...")

	err := db.AutoMigrate(
		&entity.User{},
		&entity.File{},
		&entity.Share{},    // 新增
		&entity.Favorite{}, // 新增
	)

	if err != nil {
		logger.Error("Auto migration failed", zap.Any("error", err))
		return err
	}

	logger.Info("Database auto migration completed successfully")
	return nil
}

func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
