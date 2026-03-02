package main

import (
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
	r := mux.NewRouter()
	tmpl := InitTemplates()
	r.HandleFunc("/login", userHandler.MakeHandlerLoginPage(tmpl))
	r.HandleFunc("/register", userHandler.MakeHandlerRegisterPage(tmpl))
	log.Println("Server started on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("error starting web server: %v", err)
	}
}
