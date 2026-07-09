package missions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller wires HTTP routes to the mission Service.
type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(g *echo.Group) {
	g.GET("/missions", c.ListByVehicle)
	g.POST("/missions", c.Create)
	g.POST("/missions/:id/upload", c.MarkUploaded)
}

// ListByVehicle returns all missions for a vehicle.
//
// @Summary List missions for a vehicle
// @Tags missions
// @Produce json
// @Param vehicle_id query string true "Vehicle ID"
// @Success 200 {array} missions.MissionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/missions [get]
func (c *Controller) ListByVehicle(ctx echo.Context) error {
	vehicleID := ctx.QueryParam("vehicle_id")
	if vehicleID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "vehicle_id is required"})
	}
	items, err := c.service.ListByVehicle(ctx.Request().Context(), vehicleID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, items)
}

// Create creates a new mission.
//
// @Summary Create a mission
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body missions.CreateMissionRequest true "Mission to create"
// @Success 201 {object} missions.MissionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/missions [post]
func (c *Controller) Create(ctx echo.Context) error {
	var req CreateMissionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	item, err := c.service.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, item)
}

// MarkUploaded marks a mission as uploaded to the vehicle.
//
// @Summary Mark a mission as uploaded
// @Tags missions
// @Param id path string true "Mission ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/missions/{id}/upload [post]
func (c *Controller) MarkUploaded(ctx echo.Context) error {
	if err := c.service.MarkUploaded(ctx.Request().Context(), ctx.Param("id")); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}
