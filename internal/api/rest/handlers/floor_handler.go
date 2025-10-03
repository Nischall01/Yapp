package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
	"net/http"
)

type FloorHandler struct {
	services.IFloorService
}

func NewFloorHandler(floorService services.IFloorService) *FloorHandler {
	return &FloorHandler{
		floorService,
	}
}

func (h *FloorHandler) CreateFloor(c *gin.Context) {

	u := &dto.CreateFloorReq{}

	if err := c.ShouldBindJSON(u); err != nil {

		utils.WriteError(c, utils.ErrorInvalidInput)
		return

	}

	res, err := h.IFloorService.CreateFloor(c.Request.Context(), u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
