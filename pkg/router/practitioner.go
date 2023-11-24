package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/practitioner"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/mongo"
	practService "medbuddy-backend/service/practitioner"
)

func Practitioner(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	practitionerService := practService.NewPractitionerService(dbRepo)
	practitionerCtrl := practitioner.NewController(validate, logger, practitionerService)

	practitionerUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		practitionerUrl.POST("/practitioner", practitionerCtrl.CreatePractitioner)
		practitionerUrl.GET("/practitioner", middleware.Practitioner(), practitionerCtrl.GetPractitioner)
		practitionerUrl.GET("/practitioner/email/:email", middleware.Practitioner(), practitionerCtrl.GetPractitionerByEmail)
		practitionerUrl.GET("/practitioner/ids", middleware.Practitioner(), practitionerCtrl.GetPractitionersByIds)
		practitionerUrl.GET("/practitioner/medications", middleware.Practitioner(), practitionerCtrl.GetPractitionerMedications)
	}
	return r
}
