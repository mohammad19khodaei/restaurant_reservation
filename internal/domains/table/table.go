package table

type Table struct {
	ID         int `gorm:"type:bigserial;primaryKey"`
	SeatsCount int `gorm:"type:int,NOT NULL"`
}
