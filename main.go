package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// @title                      API Recolecta
// @version                    1.0
// @description                API para gestión de recolección de residuos
// @BasePath                   /
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Escribe 'Bearer ' seguido de tu token JWT. Ejemplo: "Bearer eyJhbG..."
func main() {
	_ = godotenv.Load()

	if len(os.Args) > 1 {
		switch strings.ToLower(strings.TrimSpace(os.Args[1])) {
		case "seed", "bootstrap":
			fmt.Fprintln(os.Stderr, "bootstrap de admin deshabilitado: gestiona usuarios desde scripts/init/.env")
			os.Exit(1)
		}
	}

	InitDependencies()
}
