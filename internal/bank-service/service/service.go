package service

import (
	"NP/internal/database"
	"NP/internal/models"
	"fmt"
	"log"
)

func GetDataForWallet() (models.BankAccount, error) {
	var account models.BankAccount
	err := database.DB.First(&account).Error
	if err != nil {
		log.Printf("error get bank account: %s\n", err)
		return models.BankAccount{}, fmt.Errorf("error get bank account: %w", err)
	}
	return account, nil
}

func GetSumDonateForCategory(category string) (uint64, error) {
	var total uint64

	err := database.DB.Model(&models.Donate{}).
		Where("category = ?", category).
		Select("COALESCE(SUM(money_summ), 0)").
		Scan(&total).Error

	if err != nil {
		log.Printf("error get total donations sum for category %s: %s", category, err)
		return 0, fmt.Errorf("error get total donations sum for category %s: %w", category, err)
	}

	return total, nil
}
