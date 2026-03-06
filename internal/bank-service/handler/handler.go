package handler

import (
	"NP/internal/bank-service/service"
	"NP/internal/database"
	"NP/internal/middlewware"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MakeHandlerMyWalletPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bankAccount, err := service.GetDataForWallet()
		if err != nil {
			log.Printf("error get bank account data: %s", err)
			http.Error(w, "error get bank account data", http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "myWalletPage.html", bankAccount)
	}
}

type ForRenderDonatePage struct {
	TotalDonateForCoral   uint `json:"total_donate_for_coral"`
	TotalDonateForForests uint `json:"total_donate_for_forest"`
	TotalDonateForTurtles uint `json:"total_donate_for_turtles"`
	TotalDonateForBird    uint `json:"total_donate_for_bird"`
}

func MakeHandlerDonatePage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalSumCoral, _ := service.GetSumDonateForCategory("coral")
		totalSumForests, _ := service.GetSumDonateForCategory("forest")
		totalSumTurtles, _ := service.GetSumDonateForCategory("turtle")
		totalSumBirds, _ := service.GetSumDonateForCategory("birds")

		totalDonates := ForRenderDonatePage{
			TotalDonateForCoral:   uint(totalSumCoral),
			TotalDonateForForests: uint(totalSumForests),
			TotalDonateForTurtles: uint(totalSumTurtles),
			TotalDonateForBird:    uint(totalSumBirds),
		}
		tmpl.ExecuteTemplate(w, "donatePage.html", totalDonates)
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

	http.Redirect(w, r, "/donate", http.StatusSeeOther)

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

func MakeHandlerTopUpWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewware.GetUserIDFromRequest(r)
	if err != nil {
		log.Printf("error get user id from request: %s", err)
		return
	}
	moneySum := GetParamByUrl("moneySum", r)

	err = database.TopUpWallet(userID, uint64(moneySum))
	if err != nil {
		log.Printf("error top up wallet(user:%d money:%d): %s", userID, moneySum, err)
	}

	http.Redirect(w, r, "/my-wallet", http.StatusSeeOther)

}
