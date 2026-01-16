package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRole(allowed ...string) echo.MiddlewareFunc {
	set := map[string]bool{}
	for _, r := range allowed {
		set[r] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("role").(string)
			if role == "" || !set[role] {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
			}
			return next(c)
		}
	}
}
