package dosage

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/service/dosage"
)

type Controller struct {
	Validate      *validator.Validate
	Logger        *log.Logger
	DosageService dosage.DosageService
}

func NewController(validate *validator.Validate, logger *log.Logger, mService dosage.DosageService) *Controller {
	return &Controller{
		validate, logger, mService,
	}
}
