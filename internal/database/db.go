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
		&models.Purchases{},
		&models.Transaction{},
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
	user, err := GetUserById(userID)
	if err != nil {
		return err
	}
	bankAccount := models.BankAccount{
		Balance: 0,
		UserID:  userID,
		User:    &user,
	}
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
	}).Create(&bankAccount).Error; err != nil {
		return err
	}
	return nil
}

func TopUpWalletBalance(userID uint, amount uint64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var account models.BankAccount

		if err := tx.Where("user_id = ?", userID).First(&account).Error; err != nil {
			return fmt.Errorf("account not found: %w", err)
		}

		if err := tx.Model(&account).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		transaction := &models.Transaction{
			BankAccountID: account.ID,
			Amount:        int64(amount),
			Type:          "top_up",
			Description:   fmt.Sprintf("Пополнение счета на %d ₽", amount),
		}
		return tx.Create(transaction).Error
	})
}

func DebitWalletBalance(userID uint, amount uint64, txType string, description string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var account models.BankAccount

		if err := tx.Where("user_id = ?", userID).First(&account).Error; err != nil {
			return fmt.Errorf("account not found: %w", err)
		}

		if account.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}

		if err := tx.Model(&account).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		transaction := &models.Transaction{
			BankAccountID: account.ID,
			Amount:        -int64(amount),
			Type:          txType,
			Description:   description,
		}
		return tx.Create(transaction).Error
	})
}

func GetTransactionsHistory(userID uint) ([]models.Transaction, error) {
	var account models.BankAccount

	if err := DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	var transactions []models.Transaction
	err := DB.
		Where("bank_account_id = ?", account.ID).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, err
}

func GetTransactionsSummary(bacnkAccountID uint) (models.TransactionSummary, error) {
	var summary models.TransactionSummary

	err := DB.Model(&models.Transaction{}).
		Where("bank_account_id = ?", bacnkAccountID).
		Select(`
        COUNT(*) AS count,
        COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_top_up,
        COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END), 0) AS total_debit
    `).Scan(&summary).Error
	return summary, err
}

func PayPurchase(userID uint, amount uint64) error {
	var account models.BankAccount

	if err := DB.Where("user_id = ?", userID).First(&account).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	var cart models.Cart

	if err := DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return fmt.Errorf("cart not found: %w", err)
	}

	if account.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&account).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		purchase := models.Purchases{UserID: userID}

		if err := tx.Create(&purchase).Error; err != nil {
			return err
		}

		if err := tx.Model(&purchase).Association("PurchasesList").Append(cart.Items); err != nil {
			return err
		}

		if err := tx.Model(&cart).
			UpdateColumn("total_cost", gorm.Expr("total_cost - ?", amount)).Error; err != nil {
			return err
		}

		transaction := &models.Transaction{
			BankAccountID: account.ID,
			Amount:        -int64(amount),
			Type:          models.TxTypePurchase,
			Description:   "Purchase payment",
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return fmt.Errorf("error clearing cart items: %w", err)
		}
		return tx.Model(&cart).UpdateColumn("total_cost", 0).Error
	})
}
