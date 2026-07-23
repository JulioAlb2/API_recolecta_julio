package config

import "os"

type Config struct {
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	FCMCredentialsFile string
	ModeloReportesURL  string
	ClasificadorURL    string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		FCMCredentialsFile: os.Getenv("FCM_CREDENTIALS_FILE"),
		// Nombres de servicio de docker.compose.yml: alcanzables por DNS interno
		// de Docker dentro de app_internal_net, tanto en dev como en despliegue
		// real (mismo compose stack). Si no estan seteadas, caen a esos defaults.
		ModeloReportesURL: getEnvOrDefault("MODELO_REPORTES_URL", "http://modelo_reportes:8000"),
		ClasificadorURL:   getEnvOrDefault("CLASIFICADOR_URL", "http://clasificador_reportes:8001"),
	}
	// Aqui se pueden anadir validaciones (ej. que no esten vacios)
	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
