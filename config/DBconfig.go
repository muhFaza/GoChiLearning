package config

import (
	"fmt"
	"os"

	"gochi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func getDBConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   "gochi_db",
	}
}

func (config *DBConfig) getDBURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)
}

func InitDB() *gorm.DB {
	dbConfig := getDBConfig()
	dsn := dbConfig.getDBURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database!")
	}
	fmt.Println("Database connected!")
	return db
}
