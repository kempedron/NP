package database

import (
	"NP/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

func AddDonate(username string, sumDonate int, categoty string) error {
	var donate models.Donate

	donate.Username = username
	donate.MoneySumm = uint(sumDonate)
	donate.Category = categoty

	if err := DB.Create(&donate).Error; err != nil {
		log.Printf("error create donate: %s", err)
		return fmt.Errorf("error create donate: %w", err)
	}

	return nil

}

var DB *gorm.DB

func InitDB() error {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println("dsn: ", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}
	err = DB.AutoMigrate(
		&models.BankAccount{},
		&models.Cart{},
		&models.CartItem{},
		&models.User{},
		&models.Product{},
		&models.Donate{},
	)

	if err != nil {
		log.Printf("error migrate DB: %s", err)
		panic(err)
	}

	fmt.Println("SUCCESSFYLLU INIT DATABASE!")

	return nil

}

func GetUserById(userID uint) (models.User, error) {
	var user models.User
	err := DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return user, fmt.Errorf("error get user by id: %w", err)
	}
	return user, nil
}

func InitCart(userID uint) error {
	user, err := GetUserById(userID)
	if err != nil {
		return err
	}
	newCatr := &models.Cart{
		UserID:    userID,
		User:      &user,
		TotalCost: 0,
	}

	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(newCatr).Error; err != nil {
		return fmt.Errorf("error init cart: %w", err)
	}
	return nil
}

func SeedProducts(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		products := []models.Product{
			{Name: "Футболка 'Додо'", Price: 3200, Description: "Футболка с изображением птицы дронд"},
			{Name: "Худи 'Геккон'", Price: 5500, Description: "Худи с изображением красивого геккона"},
			{Name: "Саженец баобаба", Price: 1800, Description: "Саженец баобаба, для самостоятельного выращивания"},
			{Name: "Карта Маврикия", Price: 2900, Description: "Большая карта на стену, с отмеченными главными достопримечательностями острова"},
		}
		if err := db.Create(&products).Error; err != nil {
			return err
		}
		log.Println("Базовые товары добавлены в базу данных")
	} else {
		log.Println("Товары уже существуют, Пропуск инициализации")
	}
	return nil
}

func InitBankAccount(userID uint) error {
	bankAccount := models.BankAccount{
		Balance: 0,
		UserID:  userID,
	}
	if err := DB.Create(&bankAccount).Error; err != nil {
		return err
	}
	return nil
}

func TopUpWallet(userID uint, moneySum uint64) error {
	err := DB.Model(&models.BankAccount{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance + ?", moneySum)).Error
	if err != nil {
		log.Printf("error increment bank account balance for user %d", userID)
		return err
	}
	return nil
}
