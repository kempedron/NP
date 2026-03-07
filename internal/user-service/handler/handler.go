package handler

import (
	"NP/internal/database"
	"NP/internal/jwt"
	"NP/internal/models"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var UserAlderyExist = errors.New("пользователь уже существует")
var InvalidCredentials = errors.New("неверные учетные данные")
var SomethingWentWrong = errors.New("что-то пошло не так")

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
		fmt.Println("login func start")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
			return
		}
		if len(password) < 6 {
			http.Error(w, "Пароль должен быть не менее 6 символов", http.StatusBadRequest)
			return
		}

		var user models.User

		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		if err := user.CheckPassword(password); err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateToken(user.ID, user.Username)
		if err != nil {
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
		})

		fmt.Print("successfully set cookie2 after login")

		http.SetCookie(w, &http.Cookie{
			Name:     "logged_in",
			Value:    "true",
			HttpOnly: false,
			Path:     "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func MakeHandlerRegister(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
			return
		}
		if len(password) < 6 {
			http.Error(w, "Пароль должен быть не менее 6 символов", http.StatusConflict)
			return
		}

		var user models.User

		err := database.DB.Where("username = ?", username).First(&user).Error
		if err == nil {
			if errors.Is(err, UserAlderyExist) {
				http.Error(w, "Пользователь уже существует", http.StatusConflict)
				return
			}
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}
		newUser := &models.User{
			Username:     username,
			PasswordHash: string(hashedPassword),
		}

		if err := database.DB.Create(&newUser).Error; err != nil {
			log.Printf("ошибка создания пользователя: %s", err)
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateToken(newUser.ID, newUser.Username)
		if err != nil {
			log.Printf("ошибка создания токена: %s", err)
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		if err := database.InitCart(newUser.ID); err != nil {
			log.Printf("ошибка создания корзины: %s", err)
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		if err := database.InitBankAccount(newUser.ID); err != nil {
			log.Printf("ошибка создания банковского счета: %s", err)
			http.Error(w, "Что то пошло не так...", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "logged_in",
			Value:    "true",
			HttpOnly: false,
			Path:     "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
