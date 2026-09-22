package db

import (
	"context"
	"fmt"

	"nodistractions-online/backend/internal/domain"

	"gorm.io/gorm"
)

// schemaModels is the single registry of database tables managed by GORM.
// Add every persistent model here when the application grows.
var schemaModels = []any{
	&domain.User{},
}

// Migrate creates or updates the registered schema. Run it only through the
// migration command, which connects with DATABASE_URL_UNPOOLED.
func Migrate(ctx context.Context, gormDB *gorm.DB) error {
	if err := gormDB.WithContext(ctx).AutoMigrate(schemaModels...); err != nil {
		return fmt.Errorf("migrate schema: %w", err)
	}
	return nil
}
