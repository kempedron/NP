package main

import (
	handler2 "NP/internal/order-service/handler"
	handler "NP/internal/web-service/handler"

	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/templates/*.html"))
}

var r = mux.NewRouter()

func main() {
	tmpl := InitTemplates()
	r.HandleFunc("/", handler.MakeHandlerMainPage(tmpl))
	r.HandleFunc("/about", handler.MakeHandlerAboutPage(tmpl))
	r.HandleFunc("/buy-merch", handler.MakeHandlerBuyMerchPage(tmpl))
	r.HandleFunc("/add-to-card/{product-id:[0-9]+}/{cart-id:[0-9]+}/{quantity:[0-9]+}", handler2.AddToCart).Methods("POST")
	r.HandleFunc("/get-all-from-card", handler2.MakeGetAllCart(tmpl)).Methods("GET")
	err := http.ListenAndServe("127.0.0.1:8080", r)
	if err != nil {
		log.Fatalf("eror start web server: %v", err)
	}
}
