package handler

import (
	"NP/internal/database"
	"NP/internal/middlewware"
	"NP/internal/models"
	"NP/internal/order-service/service"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	productId := GetParamByUrl("product-id", r)
	quantity := GetParamByUrl("quantity", r)

	userID, err := middlewware.GetUserIDFromRequest(r)
	if err != nil {
		log.Printf("error get userID from jwt: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var cartID uint
	err = database.DB.Model(&models.Cart{}).
		Where("user_id = ?", userID).
		Select("id").
		First(&cartID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart := models.Cart{UserID: userID}
			if err := database.DB.Create(&cart).Error; err != nil {
				log.Printf("error creating cart for user %d: %v", userID, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			cartID = cart.ID
		} else {
			log.Printf("error get cartID from db for user %d: %s", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if err := service.AddProductByIdToCart(uint(productId), cartID, uint(quantity)); err != nil {
		log.Printf("error add product(%d) to cart: %v", productId, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Не удалось добавить товар в корзину",
		})
		return
	}

	http.Redirect(w, r, "/get-all-from-cart", http.StatusSeeOther)
}

func MakeGetAllCart(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := middlewware.GetUserIDFromRequest(r)
		if err != nil {
			log.Printf("error get userID from midlleware: %s", err)
		}
		items, err := service.GettAllProductsFromCart(uint(userID))
		if err != nil {
			log.Printf("error get products for user %d cart: %v", userID, err)
		}
		err = tmpl.ExecuteTemplate(w, "cartItems.html", items)
		if err != nil {
			log.Printf("error rendering cartItems.html: %v", err)
		}
	}
}

func GetParamByUrl(name string, r *http.Request) int {
	vars := mux.Vars(r)
	Id := vars[name]
	IdInt, err := strconv.Atoi(Id)
	if err != nil {
		log.Printf("error convert id to int: %s", err)
	}
	return IdInt
}
