package handler

import (
	"NP/internal/database"
	"NP/internal/jwt"
	"NP/internal/models"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func MakeHandlerLoginPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func MakeHandlerRegisterPage(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "registerPage.html", nil)
	}
}

func MakeHandlerLogin(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "username and password required"})
			return
		}
		if len(password) < 6 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "password must be at least 6 charactest"})
			return
		}

		var user models.User

		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			log.Printf("error check username(%s): %s", username, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid authentication"})
			return
		}
		if err := user.CheckPassword(password); err != nil {
			log.Printf("error pass check: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid authentication"})
			return
		}

		token, err := jwt.GenerateToken(user.ID, user.Username)
		if err != nil {
			log.Printf("error generate token: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Could not generate token"})
			return
		}

		cookie := new(http.Cookie)
		cookie.Name = "jwt"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.Path = "/"
		cookie.HttpOnly = true
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func MakeHandlerRegister(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "username and password required"})
			return
		}
		if len(password) < 6 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "password must be at least 6 charactest"})
			return
		}

		var user models.User

		err := database.DB.Where("username = ?", username).First(&user).Error
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "username aldery exist"})
			return
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("error db(register user): %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("error hasing password")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
		newUser := &models.User{
			Username:     username,
			PasswordHash: string(hashedPassword),
		}
		if err := database.DB.Create(&newUser).Error; err != nil {
			log.Printf("error saving user to db: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "could not create user"})
			return
		}

		token, err := jwt.GenerateToken(newUser.ID, newUser.Username)
		if err != nil {
			log.Printf("error ganerate token for new user: %v", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
