package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/pkg/handler/medicine"
	"medbuddy-backend/pkg/middleware"
	"medbuddy-backend/pkg/repository/mongo"
	medService "medbuddy-backend/service/medicine"
)

func Medicine(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *log.Logger) *gin.Engine {

	dbRepo := mongo.GetDB()
	medicineService := medService.NewMedicineService(dbRepo)
	medicineCtrl := medicine.NewController(validate, logger, medicineService)

	medicineUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		medicineUrl.POST("/medicine", middleware.Generic(), medicineCtrl.AddMedicine)
		medicineUrl.GET("/medicine/:id", middleware.Generic(), medicineCtrl.GetMedicine)
		medicineUrl.GET("/medicine", middleware.Generic(), medicineCtrl.GetMedicineFilter)
		medicineUrl.PATCH("/medicine/:id", middleware.Generic(), medicineCtrl.UpdateMedicine)
		medicineUrl.DELETE("/medicine/:id", middleware.Generic(), medicineCtrl.DeleteMedicine)
	}
	return r
}
