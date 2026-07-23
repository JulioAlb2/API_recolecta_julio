package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.GetInt("role_id")

		// El Administrador (ADMIN = 1) siempre tiene acceso a todas las rutas
		if roleID == ADMIN {
			c.Next()
			return
		}

		for _, r := range roles {
			if roleID == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "no autorizado",
		})
	}
}
