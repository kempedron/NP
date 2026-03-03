package main

import (
	bankHandler "NP/internal/bank-service/handler"
	"NP/internal/database"
	webHandler "NP/internal/web-service/handler"
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
	if err := database.InitDB(); err != nil {
		log.Fatalf("error init database for bank-service: %s", err)
	}
	r := mux.NewRouter()
	tmpl := InitTemplates()

	fmt.Println("jwt: <-", os.Getenv("JWT_SECRET"))

	r.HandleFunc("/my-wallet", bankHandler.MakeHandlerMyWalletPage(tmpl)).Methods("GET")
	r.HandleFunc("/donate", webHandler.MakeHandlerDonatePage(tmpl)).Methods("GET")
	r.HandleFunc("/donate/{category}/{username}/{moneySum}", webHandler.MakeHandlerDonate).Methods("POST")

	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
