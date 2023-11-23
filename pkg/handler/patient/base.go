package patient

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/service/patient"
)

type Controller struct {
	Validate       *validator.Validate
	Logger         *log.Logger
	PatientService patient.PatientService
}

func NewController(validate *validator.Validate, logger *log.Logger, pService patient.PatientService) *Controller {
	return &Controller{
		validate, logger, pService,
	}
}
