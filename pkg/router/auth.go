package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/patient"
	"medbuddy-backend/pkg/handler/practitioner"
	"medbuddy-backend/pkg/repository/mongo"
	patService "medbuddy-backend/service/patient"
	practService "medbuddy-backend/service/practitioner"
)

func Auth(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	patientService := patService.NewPatientService(dbRepo)
	patientCtrl := patient.NewController(validate, logger, patientService)

	practitionerService := practService.NewPractitionerService(dbRepo)
	practitionerCtrl := practitioner.Controller{Validate: validate, Logger: logger, PractitionerService: practitionerService}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/practitioner/login", practitionerCtrl.LoginPractitioner)
		authUrl.POST("/patient/login", patientCtrl.LoginPatient)
	}
	return r
}
