package user

import "time"

type User struct {
	ID        int       `gorm:"type:bigserial;primaryKey"`
	Username  string    `gorm:"type:varchar;uniqueIndex,NOT NULL"`
	Password  string    `gorm:"type:varchar,NOT NULL"`
	CreatedAt time.Time `gorm:"type:timestamp"`
}
