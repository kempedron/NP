package database

// import (
// 	"NP/internal/models"
// 	"fmt"
// 	"log"
// 	"os"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// )

// var DB *gorm.DB

// func InitDB() error {
// 	var err error
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_PORT"),
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_NAME"),
// 	)
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		log.Fatal("не удалось подключиться к БД:", err)
// 	}
// 	err = DB.AutoMigrate(
// 		&models.BankAccount{},
// 		&models.Cart{},
// 		&models.CartItem{},
// 		&models.User{},
// 		&models.Product{},
// 		&models.Donate{},
// 	)

// 	if err != nil {
// 		log.Printf("error migrate DB: %s", err)
// 		panic(err)
// 	}

// 	return nil

// }
