package main

import (
	handler "NP/internal/order-service/handler"
	"html/template"

	"github.com/gorilla/mux"
)

var tmpl *template.Template

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/add-to-card/{product-id:[0-9]+}/{cart-id:[0-9]+}/{quantity:[0-9]+}", handler.AddToCart).Methods("POST")
	r.HandleFunc("/get-all-from-card", handler.MakeGetAllCart(tmpl)).Methods("GET")

}
