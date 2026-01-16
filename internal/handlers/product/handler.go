package product

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"sales-management-api/internal/models"
	"sales-management-api/internal/services"
)

type Handler struct {
	svc *services.ProductService
	v   *validator.Validate
}

func New(svc *services.ProductService) *Handler {
	return &Handler{svc: svc, v: validator.New()}
}

type ProductReq struct {
	SKU   string `json:"sku" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,min=1"`
	Stock int64  `json:"stock" validate:"required,min=0"`
}

func (h *Handler) Create(c echo.Context) error {
	var req ProductReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	p := models.Product{
		SKU:   req.SKU,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := h.svc.Create(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, p)
}

func (h *Handler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var req ProductReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	p := models.Product{
		SKU:   req.SKU,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := h.svc.Update(uint(id), &p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "updated"})
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *Handler) List(c echo.Context) error {
	data, err := h.svc.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
