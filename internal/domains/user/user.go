package user

import "time"

type User struct {
	ID        int       `gorm:"type:bigserial;primaryKey"`
	Username  string    `gorm:"type:varchar;uniqueIndex"`
	Password  string    `gorm:"type:varchar"`
	CreatedAt time.Time `gorm:"type:timestamp"`
}
