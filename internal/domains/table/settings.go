package table

// Settings struct
type Settings struct {
	ID        int `gorm:"type:bigserial;primaryKey"`
	SeatPrice int `gorm:"type:int,NOT NULL"`
}

// TableName returns the table name
func (t Settings) TableName() string {
	return "table_settings"
}
