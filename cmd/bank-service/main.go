package main

import (
	"NP/internal/bank-service/handler"
	"NP/internal/database"
	"fmt"
	"os"

	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/templates/*.html"))
}

func main() {
	r := mux.NewRouter()
	tmpl := InitTemplates()
	if err := database.InitDB(); err != nil {
		log.Fatalf("error init db for order-service: %s", err)
	}

	fmt.Println("jwt: <-", os.Getenv("JWT_SECRET"))

	r.HandleFunc("/my-wallet", handler.MakeHandlerMyWalletPage(tmpl)).Methods("GET")
	r.HandleFunc("/donate", handler.MakeHandlerDonatePage(tmpl)).Methods("GET")
	r.HandleFunc("/donate/{category}/{username}/{moneySum}", handler.MakeHandlerDonate).Methods("POST")
	r.HandleFunc("/top-up-wallet/{moneySum}", handler.MakeHandlerTopUpWallet).Methods("POST")

	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
