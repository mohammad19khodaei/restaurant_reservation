package reservation

import "time"

type Reservation struct {
	ID         int       `gorm:"type:bigserial;primaryKey"`
	UserID     uint      `gorm:"type:int,NOT NULL"`
	TableID    uint      `gorm:"type:int,NOT NULL"`
	SeatsCount int       `gorm:"type:int,NOT NULL"`
	Price      float64   `gorm:"type:number,not null"`
	Date       time.Time `gorm:"type:timestamp,NOT NULL"`
}
