package sales

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"sales-management-api/internal/repositories"
	"sales-management-api/internal/services"
)

type Handler struct {
	svc *services.SaleService
	v   *validator.Validate
}

func New(svc *services.SaleService) *Handler {
	return &Handler{svc: svc, v: validator.New()}
}

type SaleItemReq struct {
	ProductID uint  `json:"product_id" validate:"required"`
	Qty       int64 `json:"qty" validate:"required,min=1"`
}

type CreateSaleReq struct {
	Items []SaleItemReq `json:"items" validate:"required,min=1,dive"`
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateSaleReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	// validasi minimal (validator dive untuk array)
	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// ambil cashier dari JWT middleware
	userIDAny := c.Get("user_id")
	cashierID, _ := userIDAny.(uint)

	items := make([]repositories.SaleItemInput, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, repositories.SaleItemInput{
			ProductID: it.ProductID,
			Qty:       it.Qty,
		})
	}

	sale, err := h.svc.Create(cashierID, items)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, sale)
}

func (h *Handler) List(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	limit := 50
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}

	data, err := h.svc.List(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *Handler) Detail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	sale, err := h.svc.Detail(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "sale not found"})
	}
	return c.JSON(http.StatusOK, sale)
}
