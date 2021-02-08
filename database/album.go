package database

import "gorm.io/gorm"

type Album struct {
	gorm.Model
	Title      string `json:"title"`
	OwnerRefer uint   `json:"owner_refer"`
	Owner      User   `gorm:"foreignKey:OwnerRefer" json:"owner"`
}
