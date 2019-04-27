package main

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"

	us "github.com/hellodhlyn/undersky"
)

var migration201904220001 = &gormigrate.Migration{
	ID: "201904220001",
	Migrate: func(db *gorm.DB) error {
		return db.CreateTable(
			us.User{},
		).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.DropTable(
			us.User{},
		).Error
	},
}
