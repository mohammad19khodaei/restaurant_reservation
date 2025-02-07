package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/user"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user.Repository {
	return &GormUserRepository{
		db: db,
	}
}

// Register registers a new user
func (r *GormUserRepository) Register(ctx context.Context, u *user.User) error {
	result := r.db.WithContext(ctx).Create(&u)
	if result.Error != nil {
		// handling unique_violation error
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return user.ErrUsernameAlreadyExists
		}
		return result.Error
	}

	return nil
}

// FindByUsername finds a user by username
func (r *GormUserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &u, nil
}
