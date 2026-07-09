package telemetry

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Controller wires HTTP routes to the telemetry Service.
type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(g *echo.Group) {
	g.POST("/telemetry", c.Ingest)
	g.GET("/telemetry", c.GetRange)
}

// Ingest stores a single telemetry frame.
//
// @Summary Ingest a telemetry frame
// @Tags telemetry
// @Accept json
// @Param frame body telemetry.IngestTelemetryRequest true "Telemetry frame"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/telemetry [post]
func (c *Controller) Ingest(ctx echo.Context) error {
	var req IngestTelemetryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.service.Ingest(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.NoContent(http.StatusNoContent)
}

// GetRange returns telemetry frames for a vehicle within a time range.
//
// @Summary Get telemetry frames in a time range
// @Tags telemetry
// @Produce json
// @Param vehicle_id query string true "Vehicle ID"
// @Param from query string true "Range start (RFC3339)"
// @Param to query string true "Range end (RFC3339)"
// @Success 200 {array} telemetry.TelemetryFrameResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/telemetry [get]
func (c *Controller) GetRange(ctx echo.Context) error {
	vehicleID := ctx.QueryParam("vehicle_id")
	if vehicleID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "vehicle_id is required"})
	}
	from, err := time.Parse(time.RFC3339, ctx.QueryParam("from"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid from"})
	}
	to, err := time.Parse(time.RFC3339, ctx.QueryParam("to"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid to"})
	}
	items, err := c.service.GetRange(ctx.Request().Context(), vehicleID, from, to)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, items)
}
