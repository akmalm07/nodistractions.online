package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const defaultSecretManagerSecret = "projects/933874694474/locations/us-east4/secrets/nodistractions-secrets"

type Config struct {
	ListenAddress string
	FrontendURL   string
	DatabaseURL   string
	SessionSecret string
	JWTIssuer     string
	JWTAudience   string
}

// Load reads the API configuration. Local development reads .env; secret mode
// reads a dotenv-formatted payload from Google Secret Manager instead.
func Load(ctx context.Context) (Config, error) {
	environment, err := loadEnvironment(ctx)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		ListenAddress: ":" + value(environment, "PORT", "3000"),
		FrontendURL:   value(environment, "FRONTEND_URL", "http://localhost:5173"),
		DatabaseURL:   value(environment, "DATABASE_URL", ""),
		SessionSecret: value(environment, "SESSION_SECRET", ""),
		JWTIssuer:     value(environment, "JWT_ISSUER", ""),
		JWTAudience:   value(environment, "JWT_AUDIENCE", ""),
	}

	for name, value := range map[string]string{
		"DATABASE_URL":   config.DatabaseURL,
		"SESSION_SECRET": config.SessionSecret,
		"JWT_ISSUER":     config.JWTIssuer,
		"JWT_AUDIENCE":   config.JWTAudience,
	} {
		if value == "" {
			return Config{}, fmt.Errorf("missing required environment variable: %s", name)
		}
	}

	return config, nil
}

// MigrationDatabaseURL returns the direct Neon URL, never the pooled URL.
func MigrationDatabaseURL(ctx context.Context) (string, error) {
	environment, err := loadEnvironment(ctx)
	if err != nil {
		return "", err
	}

	url := value(environment, "DATABASE_URL_UNPOOLED", "")
	if url == "" {
		return "", fmt.Errorf("missing required environment variable: DATABASE_URL_UNPOOLED")
	}
	return url, nil
}

func loadDotEnv() {
	// The first path supports `cd backend`; the second supports a backend-local .env.
	_ = godotenv.Load("../.env", ".env")
}

func loadEnvironment(ctx context.Context) (map[string]string, error) {
	secretMode, err := strconv.ParseBool(valueFromOS("ENV_MODE", "false"))
	if err != nil {
		return nil, fmt.Errorf("ENV_MODE must be true or false: %w", err)
	}

	if !secretMode {
		loadDotEnv()
		return processEnvironment(), nil
	}

	environment := processEnvironment()
	secretName := value(environment, "SECRET_MANAGER_SECRET", defaultSecretManagerSecret)
	secretEnvironment, err := loadSecretEnvironment(ctx, secretName)
	if err != nil {
		return nil, err
	}
	for key, value := range secretEnvironment {
		environment[key] = value
	}
	return environment, nil
}

func loadSecretEnvironment(ctx context.Context, secretName string) (map[string]string, error) {
	versionName, err := secretVersionName(secretName)
	if err != nil {
		return nil, err
	}

	clientOptions := []option.ClientOption{}
	if location := secretLocation(versionName); location != "" {
		clientOptions = append(clientOptions, option.WithEndpoint(
			fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", location),
		))
	}

	client, err := secretmanager.NewClient(ctx, clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("create Secret Manager client: %w", err)
	}
	defer client.Close()

	response, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionName,
	})
	if err != nil {
		return nil, fmt.Errorf("access application configuration secret: %w", err)
	}

	environment, err := godotenv.Unmarshal(string(response.Payload.Data))
	if err != nil {
		return nil, fmt.Errorf("parse application configuration secret as dotenv: %w", err)
	}
	return environment, nil
}

func secretVersionName(secretName string) (string, error) {
	name := strings.TrimSuffix(strings.TrimSpace(secretName), "/")
	if !strings.HasPrefix(name, "projects/") {
		return "", fmt.Errorf("SECRET_MANAGER_SECRET must be a full projects/... secret resource name")
	}
	if strings.Contains(name, "/versions/") {
		return name, nil
	}
	return name + "/versions/latest", nil
}

func secretLocation(secretName string) string {
	segments := strings.Split(secretName, "/")
	for index, segment := range segments {
		if segment == "locations" && index+1 < len(segments) {
			return segments[index+1]
		}
	}
	return ""
}

func processEnvironment() map[string]string {
	environment := make(map[string]string)
	for _, pair := range os.Environ() {
		key, value, found := strings.Cut(pair, "=")
		if found {
			environment[key] = value
		}
	}
	return environment
}

func valueFromOS(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func value(environment map[string]string, key, fallback string) string {
	if value := environment[key]; value != "" {
		return value
	}
	return fallback
}
