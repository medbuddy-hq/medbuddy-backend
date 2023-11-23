package practitioner

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"medbuddy-backend/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *log.Logger
}

func (base *Controller) LoginPractitioner(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "practitioner created successfully", gin.H{"practitioner": "practitioner object"})
	c.JSON(http.StatusOK, rd)

}
