package userlogin

import (
	"github.com/geoikonomou/forum/internal/app"
	"github.com/geoikonomou/forum/internal/config"
	"github.com/geoikonomou/forum/internal/domain/session"
	"github.com/geoikonomou/forum/internal/infra/logger"
)

type Handler struct {
	UserServices   app.Services
	SessionManager session.Manager
	Config         *config.ServerConfig
	Logger         logger.Logger
}

func NewHandler(config *config.ServerConfig, app app.Services, sm session.Manager, logger logger.Logger) *Handler {
	return &Handler{
		UserServices:   app,
		SessionManager: sm,
		Config:         config,
		Logger:         logger,
	}
}
