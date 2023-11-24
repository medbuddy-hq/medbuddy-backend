package practitioner

import (
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"medbuddy-backend/utility"
)

func (base *Controller) CreatePractitioner(c *gin.Context) {
	var data model.PractitionerRequest

	if err := c.BindJSON(&data); err != nil {
		base.Logger.Error("Error when binding request body, error: ", err.Error())
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrRequest, constant.ErrRequest, nil)
		c.JSON(rd.Code, rd)
		return
	}

	if err := base.Validate.Struct(data); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrValidation, err.Error(), nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.PractitionerService.CreatePractitioner(&data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) LoginPractitioner(c *gin.Context) {

	var data model.UserLogin

	if err := c.BindJSON(&data); err != nil {
		base.Logger.Error("Error when binding request body on LoginPractitioner, error: ", err.Error())
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrRequest, constant.ErrRequest, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validate.Struct(data); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrValidation, err.Error(), nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.PractitionerService.LoginPractitioner(&data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, http.StatusText(err.Code()), err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPractitioner(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	userInfo := uInfo.(*model.ContextInfo)
	response, err := base.PractitionerService.GetPractitioner(userInfo)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPractitionerByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrValidation, constant.ErrRequest, nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.PractitionerService.GetPractitionerByEmail(email)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPractitionersByIds(c *gin.Context) {
	var data []string

	if err := c.BindJSON(&data); err != nil {
		base.Logger.Error("Error when binding request body, error: ", err.Error())
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrRequest, constant.ErrRequest, nil)
		c.JSON(rd.Code, rd)
		return
	}

	if len(data) <= 0 {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "no practitioner id's specified", nil)
		c.JSON(rd.Code, rd)
		return
	}

	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userInfo := uInfo.(*model.ContextInfo)

	response, err := base.PractitionerService.GetPractitionersByIDs(userInfo, data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPractitionerMedications(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userInfo := uInfo.(*model.ContextInfo)

	response, err := base.PractitionerService.GetPractitionerMedications(userInfo)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}
