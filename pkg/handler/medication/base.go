package medication

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/service/medication"
)

type Controller struct {
	Validate          *validator.Validate
	Logger            *log.Logger
	MedicationService medication.MedicationService
}

func NewController(validate *validator.Validate, logger *log.Logger, mService medication.MedicationService) *Controller {
	return &Controller{
		validate, logger, mService,
	}
}
