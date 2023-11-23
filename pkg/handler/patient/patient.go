package patient

import (
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"medbuddy-backend/utility"
)

func (base *Controller) LoginPatient(c *gin.Context) {
	var data model.UserLogin

	if err := c.BindJSON(&data); err != nil {
		base.Logger.Error("Error when binding request body on LoginPatient, error: ", err.Error())
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrRequest, constant.ErrRequest, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validate.Struct(data); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrValidation, err.Error(), nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.PatientService.LoginPatient(&data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, http.StatusText(err.Code()), err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "patient logged in successfully", response)
	c.JSON(rd.Code, rd)

}

func (base *Controller) CreatePatient(c *gin.Context) {
	var data model.CreatePatientReq

	if err := c.Bind(&data); err != nil {
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

	response, err := base.PatientService.CreatePatient(&data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetPatient(c *gin.Context) {
	uInfo, exists := c.Get("user info")
	if !exists {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, constant.StatusFailed, constant.ErrServer, constant.ErrRequest, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	uId := uInfo.(*model.ContextInfo).ID
	response, err := base.PatientService.GetPatient(uId)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}
