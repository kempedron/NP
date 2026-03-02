package main

import (
	webHandler "NP/internal/web-service/handler"
	"NP/internal/web-service/service"

	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitTemplates() *template.Template {
	return template.Must(template.ParseGlob("web/templates/*.html"))
}

func main() {
	if err := service.InitDB(); err != nil {
		log.Fatalf("error init database: %s", err)
	}
	r := mux.NewRouter()
	tmpl := InitTemplates()

	r.HandleFunc("/", webHandler.MakeHandlerMainPage(tmpl)).Methods("GET")
	r.HandleFunc("/about-us", webHandler.MakeHandlerAboutPage(tmpl)).Methods("GET")
	r.HandleFunc("/buy-merch", webHandler.MakeHandlerBuyMerchPage(tmpl)).Methods("GET")
	r.HandleFunc("/donate", webHandler.MakeHandlerDonatePage(tmpl)).Methods("GET")
	r.HandleFunc("/donate/{category}/{username}/{moneySum}", webHandler.MakeHandlerDonate).Methods("POST")

	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
