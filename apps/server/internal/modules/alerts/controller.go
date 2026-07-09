package alerts

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller wires HTTP routes to the alert Service.
type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(g *echo.Group) {
	g.POST("/alerts", c.Create)
	g.GET("/alerts", c.ListRecent)
}

// Create records a new alert for a vehicle.
//
// @Summary Create an alert
// @Tags alerts
// @Accept json
// @Param alert body alerts.CreateAlertRequest true "Alert to create"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/alerts [post]
func (c *Controller) Create(ctx echo.Context) error {
	var req CreateAlertRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.service.Create(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}

// ListRecent returns the most recent alerts for a vehicle.
//
// @Summary List recent alerts for a vehicle
// @Tags alerts
// @Produce json
// @Param vehicle_id query string true "Vehicle ID"
// @Success 200 {array} alerts.AlertResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/alerts [get]
func (c *Controller) ListRecent(ctx echo.Context) error {
	vehicleID := ctx.QueryParam("vehicle_id")
	if vehicleID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "vehicle_id is required"})
	}
	items, err := c.service.ListRecent(ctx.Request().Context(), vehicleID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, items)
}
