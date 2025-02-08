package repositories

import (
	"context"

	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/table"
	"gorm.io/gorm"
)

// GormTableRepository struct
type GormTableRepository struct {
	db *gorm.DB
}

// NewGormTableRepository creates a new GormTableRepository
func NewGormTableRepository(db *gorm.DB) table.Repository {
	return &GormTableRepository{db: db}
}

// CreateTable creates a new table
func (r *GormTableRepository) CreateTable(ctx context.Context, table *table.Table) error {
	return r.db.WithContext(ctx).Create(table).Error
}

// GetTotalCount returns the total number of tables
func (r *GormTableRepository) GetTotalCount(ctx context.Context) int {
	var count int64
	r.db.WithContext(ctx).Model(&table.Table{}).Count(&count)
	return int(count)
}

// CreateTableSettings creates a new table settings
func (r *GormTableRepository) CreateTableSettings(ctx context.Context, seatPrice int) error {
	return r.db.WithContext(ctx).Create(&table.Settings{SeatPrice: seatPrice}).Error
}
