package practitioner

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/service/practitioner"
)

type Controller struct {
	Validate            *validator.Validate
	Logger              *log.Logger
	PractitionerService practitioner.PractitionerService
}

func NewController(validate *validator.Validate, logger *log.Logger, pService practitioner.PractitionerService) *Controller {
	return &Controller{
		validate, logger, pService,
	}
}
