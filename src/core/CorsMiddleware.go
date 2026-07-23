package core

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsEnv == "" {
		originsEnv = os.Getenv("CORS_ORIGINS")
	}

	// En desarrollo, debug o si es wildcard, permitimos todos los orígenes dinámicamente
	isDev := os.Getenv("ENVIRONMENT") == "development" || os.Getenv("DEBUG") == "true"

	var origins []string
	hasWildcard := false
	for _, origin := range strings.Split(originsEnv, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			if origin == "*" {
				hasWildcard = true
			}
			origins = append(origins, origin)
		}
	}

	if len(origins) == 0 || hasWildcard || isDev {
		return cors.New(cors.Config{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

