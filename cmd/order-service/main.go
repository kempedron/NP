package main

import (
	"html/template"
	"log"
	"net/http"

	"NP/internal/database"
	"NP/internal/order-service/handler"

	"github.com/gorilla/mux"
)

func InitTemplates() *template.Template {
	funcMap := template.FuncMap{
		"multiply": func(a, b uint) uint {
			return a * b
		},
	}
	return template.Must(
		template.New("").Funcs(funcMap).ParseGlob("web/templates/*.html"),
	)
}

func main() {
	r := mux.NewRouter()
	tmpl := InitTemplates()

	if err := database.InitDB(); err != nil {
		log.Fatalf("error init db for order-service: %s", err)
	}

	if err := database.SeedProducts(database.DB); err != nil {
		log.Fatalf("error seed db for order-service: %s", err)
	}

	r.HandleFunc("/add-to-cart/{product-id:[0-9]+}/{quantity:[0-9]+}", handler.AddToCart).Methods("POST")
	r.HandleFunc("/get-all-from-cart", handler.MakeGetAllCart(tmpl)).Methods("GET")
	r.HandleFunc("/purchase-cart", handler.MakePurchaseCart(tmpl)).Methods("POST")
	r.HandleFunc("/my-purchases", handler.MakeGetMyPurchases(tmpl)).Methods("GET")
	r.HandleFunc("/delete-from-cart/{product-id:[0-9]+}", handler.DeletePurchaseFromCartHandler).Methods("POST")

	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
