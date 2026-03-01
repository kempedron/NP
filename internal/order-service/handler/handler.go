package handler

import (
	"NP/internal/order-service/service"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {
	productId := GetParamByUrl("product-id", r)
	cartId := GetParamByUrl("cart-id", r)
	quantity := GetParamByUrl("quantity", r)

	if err := service.AddProductByIdToCart(uint(productId), uint(cartId), uint(quantity)); err != nil {
		log.Printf("error add product(%d) to cart: %v", productId, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Не удалось добавить товар в корзину",
		})
	}

}

func MakeGetAllCart(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cartId := GetParamByUrl("cart-id", r)
		err, items := service.GettAllProductsFromCart(uint(cartId))
		if err != nil {
			log.Printf("error get products for cart %d: %v", cartId, err)
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
