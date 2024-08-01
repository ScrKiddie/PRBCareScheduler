package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewDatabase(config *viper.Viper) *gorm.DB {
	username := config.GetString("db.username")
	password := config.GetString("db.password")
	host := config.GetString("db.host")
	port := config.GetInt("db.port")
	database := config.GetString("db.name")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&statement_timeout=5000", username, password, host, port, database)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
