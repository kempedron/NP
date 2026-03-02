package service

import (
	"NP/internal/models"
	"NP/internal/web-service/service"
)

var DB = service.DB

func AddUserToDB(username string, password string) error {
	var user models.User
	user.Username = username
	user.HashPassword(password)
	return nil
}
