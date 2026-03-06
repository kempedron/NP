package service

import (
	"NP/internal/database"
	"NP/internal/models"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddProductByIdToCart(productId uint, cartId uint, quantity uint) error {

	catrItem := models.CartItem{
		CartID:    cartId,
		ProductID: productId,
		Quantity:  quantity,
	}

	err := database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "cart_id"}, {Name: "product_id"}},
		DoUpdates: clause.Set{
			{Column: clause.Column{Name: "quantity"}, Value: gorm.Expr("carts_item.quantity + ?", quantity)},
		},
	}).Create(&catrItem).Error

	if err != nil {
		log.Printf("error add/upd product in cart: %s", err)
		return fmt.Errorf("error add/upd product in cart: %w", err)
	}
	return nil

}

type RespForGetProducts struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       uint   `json:"price"`
	Quantity    uint   `json:"quantity"`
	Description string `json:"description"`
}

func GettAllProductsFromCart(userID uint) ([]RespForGetProducts, error) {
	var cart models.Cart

	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []RespForGetProducts{}, nil
		}
		log.Printf("error find cart for user %d: %s", userID, err)
		return nil, fmt.Errorf("error find cart for user %d: %w", userID, err)
	}
	var items []models.CartItem

	err := database.DB.
		Preload("Product").
		Where("cart_id = ?", cart.ID).
		Find(&items).Error
	if err != nil {
		log.Printf("error get cart items for user %d: %v", userID, err)
		return nil, fmt.Errorf("error get cart items for user %d: %w", userID, err)
	}

	var result []RespForGetProducts

	for ind, val := range items {
		n := RespForGetProducts{
			ID:          ind,
			Name:        val.Product.Name,
			Price:       val.Product.Price,
			Quantity:    val.Quantity,
			Description: val.Product.Description,
		}

		result = append(result, n)
	}

	return result, nil
}
