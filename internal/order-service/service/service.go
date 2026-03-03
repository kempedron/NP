package service

import (
	"NP/internal/database"
	"NP/internal/models"
	"fmt"
	"log"
)

var DB = database.DB

func AddProductByIdToCart(productId uint, cartId uint, quantity uint) error {

	catrItem := models.CartItem{
		CartID:    cartId,
		ProductID: productId,
		Quantity:  quantity,
	}

	if err := DB.Create(&catrItem).Error; err != nil {
		log.Printf("error add product to cart: %v", err)
		return fmt.Errorf("error add product to cart: %w", err)
	}
	return nil

}

func GettAllProductsFromCart(cartID uint) (error, []models.CartItem) {
	var cartItems []models.CartItem

	err := DB.
		Where("cart_id = ?", cartID).
		Preload("Cart").
		Preload("Product").
		Find(&cartItems).Error
	if err != nil {
		log.Printf("error get cart(%d): %v", cartID, err)
		return fmt.Errorf("error get cart(%d): %w", cartID, err), nil
	}

	return nil, cartItems
}
