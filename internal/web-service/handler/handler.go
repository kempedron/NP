package handler

import (
	"NP/internal/database"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MakeHandlerMainPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	}
}

func MakeHandlerAboutPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "about.html", nil)
	}
}

func MakeHandlerBuyMerchPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "buyMerch.html", nil)
	}
}

func MakeHandlerDonatePage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "donatePage.html", nil)
	}
}

func MakeHandlerDonate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]
	username := vars["username"]
	moneySum := GetParamByUrl("moneySum", r)

	err := database.AddDonate(username, moneySum, category)
	if err != nil {
		log.Printf("error donating(user:%s category:%s money:%d): %s", username, category, moneySum, err)
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

// donate/{category}/{username}/{moneySum}
