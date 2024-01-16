package models

import "gorm.io/gorm"

type Pharmacy struct {
	ID      uint    `gorm:"primary key; autoincrement" json:"id"`
	Owner   *string `json:"owner"`
	Name    *string `json:"name"`
	Address *string `json:"address"`
}

func MigratePharmacy(db *gorm.DB) error {
	err := db.AutoMigrate(&Pharmacy{})
	return err
}
