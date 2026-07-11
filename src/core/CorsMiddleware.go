package core

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	originsEnv := os.Getenv("CORS_ORIGINS")
	if originsEnv == "" {
		originsEnv = os.Getenv("CORS_ALLOWED_ORIGINS")
	}

	var origins []string
	for _, origin := range strings.Split(originsEnv, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	// Si no hay orígenes configurados, permitir todos (útil en dev/staging sin variable seteada).
	allowHeaders := []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"ngrok-skip-browser-warning",
	}

	if len(origins) == 0 {
		return cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:    allowHeaders,
			MaxAge:          12 * time.Hour,
		})
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

