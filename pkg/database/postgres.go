package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wandermaia/customer-api/internal/config"
	"github.com/wandermaia/customer-api/internal/domain/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresConnection cria uma nova conexão com o banco de dados PostgreSQL
func NewPostgresConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Auto-migra as tabelas
	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		return nil, err
	}

	return db, nil
}
