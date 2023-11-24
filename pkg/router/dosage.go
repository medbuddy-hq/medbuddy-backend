package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/dosage"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/mongo"
	dosService "medbuddy-backend/service/dosage"
)

func Dosage(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	dosageService := dosService.NewDosageService(dbRepo)
	dosageCtrl := dosage.NewController(validate, logger, dosageService)

	dosageUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		dosageUrl.PATCH("/dosage-status/:id", middleware.Patient(), dosageCtrl.UpdateDosageStatus)
		dosageUrl.GET("/dosage/:id", middleware.Patient(), dosageCtrl.GetDosage)
	}
	return r
}
