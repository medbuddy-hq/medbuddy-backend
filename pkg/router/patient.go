package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/patient"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/mongo"
	patService "medbuddy-backend/service/patient"
)

func Patient(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	patientService := patService.NewPatientService(dbRepo)
	patientCtrl := patient.NewController(validate, logger, patientService)

	patientUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		patientUrl.POST("/patient", patientCtrl.CreatePatient)
		patientUrl.GET("/patient", middleware.Patient(), patientCtrl.GetPatient)
		//patientUrl.GET("/patient/:id", patientCtrl.GetPatientByID)
		//patientUrl.PATCH("/patient", patientCtrl.UpdatePatient)
	}
	return r
}
