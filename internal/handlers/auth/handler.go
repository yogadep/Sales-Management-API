package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"sales-management-api/internal/models"
	"sales-management-api/internal/services"
)

type Handler struct {
	svc *services.AuthService
	v   *validator.Validate
}

func New(svc *services.AuthService) *Handler {
	return &Handler{svc: svc, v: validator.New()}
}

type LoginReq struct {
	Username string `json:"username" validate:"required,min=3,max=80"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, role, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"role":  role,
	})
}

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=80"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	Role     string `json:"role" validate:"required,oneof=admin kasir"`
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := h.v.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.svc.Register(req.Username, req.Password, models.Role(req.Role)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "registered"})
}
