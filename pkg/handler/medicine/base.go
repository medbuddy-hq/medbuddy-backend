package medicine

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/service/medicine"
)

type Controller struct {
	Validate        *validator.Validate
	Logger          *log.Logger
	MedicineService medicine.MedicineService
}

func NewController(validate *validator.Validate, logger *log.Logger, mService medicine.MedicineService) *Controller {
	return &Controller{
		validate, logger, mService,
	}
}
