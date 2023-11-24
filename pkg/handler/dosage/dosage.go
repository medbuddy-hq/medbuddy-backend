package dosage

import (
	"github.com/gin-gonic/gin"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/utility"
	"net/http"
	"strings"
)

func (base *Controller) GetPatientDosages(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	var isActive *bool
	userInfo := uInfo.(*model.ContextInfo)
	isActiveParam, exists := c.GetQuery("is_active")
	if exists {
		isActiveTemp := strings.TrimSpace(strings.ToLower(isActiveParam)) == "true"
		isActive = &isActiveTemp
	}

	dosages, err := base.DosageService.GetPatientsDosages(*userInfo, isActive, "")
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "", dosages)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetMedicationDosages(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	var isActive *bool
	userInfo := uInfo.(*model.ContextInfo)
	medId, _ := c.GetQuery("medication_id")
	isActiveParam, exists := c.GetQuery("is_active")
	if exists {
		isActiveTemp := strings.TrimSpace(strings.ToLower(isActiveParam)) == "true"
		isActive = &isActiveTemp
	}

	dosages, err := base.DosageService.GetPatientsDosages(*userInfo, isActive, medId)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "", dosages)
	c.JSON(rd.Code, rd)
}

func (base *Controller) UpdateDosageStatus(c *gin.Context) {
	var data model.SetStatusRequest

	uInfo, exists := c.Get("user info")
	if !exists {
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

	dosageId := c.Param("id")
	if err := base.DosageService.SetDosageStatus(userInfo, data.Status, dosageId); err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "dosage status updated successfully", nil)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetDosage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.DosageService.GetDosage(id)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}
