package core

import (
	"log"
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

	allowHeaders := []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"ngrok-skip-browser-warning",
	}
	allowMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

	env := strings.ToLower(strings.TrimSpace(os.Getenv("ENVIRONMENT")))
	isProd := env == "production" || env == "prod"

	// Seguridad: nunca abrir CORS a todos en producción.
	// Si no hay orígenes configurados en dev, AllowAllOrigins (útil sin variable seteada).
	if len(origins) == 0 {
		if isProd {
			log.Println("CORS: ENVIRONMENT=production sin CORS_ORIGINS — se deniegan orígenes cruzados")
			return cors.New(cors.Config{
				AllowOrigins:     []string{},
				AllowMethods:     allowMethods,
				AllowHeaders:     allowHeaders,
				AllowCredentials: false,
				MaxAge:           12 * time.Hour,
			})
		}
		log.Println("CORS: sin lista de orígenes en desarrollo — AllowAllOrigins=true (solo local/dev)")
		return cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    allowMethods,
			AllowHeaders:    allowHeaders,
			MaxAge:          12 * time.Hour,
		})
	}

	log.Printf("CORS: orígenes permitidos = %v", origins)
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
