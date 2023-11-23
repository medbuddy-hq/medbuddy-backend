package medicine

import (
	"github.com/gin-gonic/gin"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/utility"
	"net/http"
)

func (base *Controller) AddMedicine(c *gin.Context) {
	var data model.MedicineRequest

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

	if err := base.MedicineService.AddMedicine(&data); err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "successfully added medicine", nil)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetMedicine(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.MedicineService.GetMedicine(id)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) GetMedicineFilter(c *gin.Context) {
	name, _ := c.GetQuery("name")
	manufacturer, _ := c.GetQuery("manufacturer")
	form, _ := c.GetQuery("form")
	strength, _ := c.GetQuery("strength")

	medFilter := model.MedicineFilter{
		Name:         name,
		Manufacturer: manufacturer,
		Form:         form,
		Strength:     strength,
	}

	if err := base.Validate.Struct(medFilter); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrValidation, "ensure all required query parameters are set", nil)
		c.JSON(rd.Code, rd)
		return
	}

	response, err := base.MedicineService.GetMedicineFilter(&medFilter)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) UpdateMedicine(c *gin.Context) {
	var data model.MedicineRequest
	id := c.Param("id")

	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

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

	response, err := base.MedicineService.UpdateMedicine(id, &data)
	if err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "medicine updated successfully", response)
	c.JSON(rd.Code, rd)
}

func (base *Controller) DeleteMedicine(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, constant.StatusFailed, constant.ErrRequest, "missing id parameter", nil)
		c.JSON(rd.Code, rd)
		return
	}

	if err := base.MedicineService.DeleteMedicine(id); err != nil {
		rd := utility.BuildErrorResponse(err.Code(), constant.StatusFailed, constant.ErrRequest, err.Error(), nil)
		c.JSON(err.Code(), rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "medicine deleted successfully", nil)
	c.JSON(rd.Code, rd)
}
