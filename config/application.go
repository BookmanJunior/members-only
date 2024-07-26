package config

import (
	"log"

	"github.com/bookmanjunior/members-only/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Users    *models.UserModel
	Messages *models.MessageModel
	Avatar   *models.AvatarModel
}
