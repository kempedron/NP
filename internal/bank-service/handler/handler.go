package handler

import (
	"NP/internal/bank-service/service"
	"html/template"
	"net/http"
)

func MakeHandlerMyWalletPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := service.GetDataForWallet()
		tmpl.ExecuteTemplate(w, "myWalletPage.html", data)
	}
}
