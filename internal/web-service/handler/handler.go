package handler

import (
	"html/template"
	"net/http"
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
