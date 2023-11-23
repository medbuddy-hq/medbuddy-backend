package router

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"medbuddy-backend/pkg/handler/health"
)

func Health(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	health := health.Controller{Validate: validate, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/health", health.Post)
		authUrl.GET("/health", health.Get)
	}
	return r
}
