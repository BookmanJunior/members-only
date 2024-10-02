package config

import (
	"log"

	"github.com/bookmanjunior/members-only/internal/cloud"
	"github.com/bookmanjunior/members-only/internal/hub"
	"github.com/bookmanjunior/members-only/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Users    *models.UserModel
	Messages *models.MessageModel
	Avatar   *models.AvatarModel
	Servers  *models.ServerModel
	Channels *models.ChannelModel
	Cloud    cloud.Cloud
	Hub      *hub.Hub
}
