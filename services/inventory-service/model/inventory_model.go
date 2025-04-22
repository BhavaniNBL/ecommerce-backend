package model

type Inventory struct {
	ProductID         string `gorm:"primaryKey"`
	Quantity          int32
	WarehouseLocation string
}
