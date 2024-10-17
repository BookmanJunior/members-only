package config

import (
	"log"

	"github.com/bookmanjunior/members-only/internal/cloud"
	"github.com/bookmanjunior/members-only/internal/hub"
	"github.com/bookmanjunior/members-only/internal/models"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Users          *models.UserModel
	Messages       *models.MessageModel
	Avatar         *models.AvatarModel
	Servers        *models.ServerModel
	ServerMembers  *models.ServerMembersModel
	ServerMessages *models.ServerMessageModel
	Channels       *models.ChannelModel
	Redis          *redis.Client
	Cloud          cloud.Cloud
	Hub            *hub.Hub
}
