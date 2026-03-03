package main

import (
	"NP/internal/database"
	userHandler "NP/internal/user-service/handler"
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
		log.Fatalf("error init db for user-service: %s", err)
	}
	r := mux.NewRouter()
	tmpl := InitTemplates()
	r.HandleFunc("/login", userHandler.MakeHandlerLoginPage(tmpl)).Methods("GET")
	r.HandleFunc("/register", userHandler.MakeHandlerRegisterPage(tmpl)).Methods("GET")

	r.HandleFunc("/login", userHandler.MakeHandlerLogin(tmpl)).Methods("POST")
	r.HandleFunc("/register", userHandler.MakeHandlerRegister(tmpl)).Methods("POST")

	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
