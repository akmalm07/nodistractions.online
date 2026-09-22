package db

import (
	"context"
	"errors"

	"nodistractions-online/backend/internal/domain"

	"gorm.io/gorm"
)

// UserRepository owns user persistence operations.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(gormDB *gorm.DB) *UserRepository {
	return &UserRepository{db: gormDB}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrEmailTaken
	}
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, ErrNotFound
	}
	return user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, ErrNotFound
	}
	return user, err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
