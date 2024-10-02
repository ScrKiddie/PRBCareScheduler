package config

import (
	"fmt"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewDatabase(config *viper.Viper) *gorm.DB {
	username := config.GetString("db.username")
	password := config.GetString("db.password")
	host := config.GetString("db.host")
	port := config.GetInt("db.port")
	database := config.GetString("db.name")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&lock_timeout=5000", username, password, host, port, database)

	logger := slogGorm.New(
		slogGorm.WithRecordNotFoundError(),
		slogGorm.WithSlowThreshold(500*time.Millisecond),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
