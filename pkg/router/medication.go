package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/dosage"
	"medbuddy-backend/pkg/handler/medication"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/mongo"
	dosService "medbuddy-backend/service/dosage"
	medService "medbuddy-backend/service/medication"
)

func Medication(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	medicationService := medService.NewMedicationService(dbRepo)
	medicationCtrl := medication.NewController(validate, logger, medicationService)
	dosageService := dosService.NewDosageService(dbRepo)
	dosageCtrl := dosage.NewController(validate, logger, dosageService)

	medicationUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		medicationUrl.POST("/medication", middleware.Patient(), medicationCtrl.AddMedication)
		medicationUrl.GET("/medication/:id", middleware.Patient(), medicationCtrl.GetMedication)
		medicationUrl.GET("/medication/dosages", middleware.Patient(), dosageCtrl.GetMedicationDosages)
		medicationUrl.GET("/medication", middleware.Patient(), medicationCtrl.GetPatientMedications)
		medicationUrl.DELETE("/medication/:id", middleware.Patient(), medicationCtrl.DeleteMedication)
	}
	return r
}
