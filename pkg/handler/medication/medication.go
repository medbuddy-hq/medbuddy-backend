package medication

import (
	"github.com/gin-gonic/gin"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/utility"
	"net/http"
)

func (base *Controller) AddMedication(c *gin.Context) {
	var data model.MedicationRequest

	uInfo, exists := c.Get("user info")
	if !exists {
		base.Logger.Error("could not find patient's details in token")
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userInfo := uInfo.(*model.ContextInfo)

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

	res, err := base.MedicationService.AddMedication(userInfo, &data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successfully added medication", res)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetMedication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.MedicationService.GetMedication(id)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPatientMedications(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}
	userInfo := uInfo.(*model.ContextInfo)

	response, err := base.MedicationService.GetPatientMedications(*userInfo)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) DeleteMedication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
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

	if err := base.MedicationService.DeleteMedication(userInfo, id); err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "deleted medication successfully", nil)
	c.JSON(rd.Code, rd)
}

func (base *Controller) AddPractitionerToMeds(c *gin.Context) {
	id := c.Param("medication-id")
	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

	var emails []string
	if err := c.BindJSON(&emails); err != nil {
		base.Logger.Error("Error when binding request body, error: ", err.Error())
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrRequest, constant.ErrRequest, nil)
		c.JSON(rd.Code, rd)
		return
	}

	if len(emails) <= 0 {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "no practitioner email specified", nil)
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

	response, err := base.MedicationService.AddPractitionersToMedication(userInfo, id, emails)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, response, nil)
	c.JSON(rd.Code, rd)
}
