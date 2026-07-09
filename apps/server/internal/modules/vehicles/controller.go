package vehicles

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller wires HTTP routes to the vehicle Service.
type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterRoutes(g *echo.Group) {
	g.GET("/vehicles", c.List)
	g.GET("/vehicles/:id", c.Get)
	g.POST("/vehicles", c.Create)
}

// List returns all vehicles.
//
// @Summary List vehicles
// @Tags vehicles
// @Produce json
// @Success 200 {array} vehicles.VehicleResponse
// @Failure 500 {object} map[string]string
// @Router /api/vehicles [get]
func (c *Controller) List(ctx echo.Context) error {
	items, err := c.service.List(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, items)
}

// Get returns a vehicle by ID.
//
// @Summary Get a vehicle
// @Tags vehicles
// @Produce json
// @Param id path string true "Vehicle ID"
// @Success 200 {object} vehicles.VehicleResponse
// @Failure 404 {object} map[string]string
// @Router /api/vehicles/{id} [get]
func (c *Controller) Get(ctx echo.Context) error {
	item, err := c.service.Get(ctx.Request().Context(), ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "vehicle not found"})
	}
	return ctx.JSON(http.StatusOK, item)
}

// Create creates a new vehicle.
//
// @Summary Create a vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body vehicles.CreateVehicleRequest true "Vehicle to create"
// @Success 201 {object} vehicles.VehicleResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/vehicles [post]
func (c *Controller) Create(ctx echo.Context) error {
	var req CreateVehicleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	item, err := c.service.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, item)
}
