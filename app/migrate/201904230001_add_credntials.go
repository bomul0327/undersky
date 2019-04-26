package main

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"

	us "github.com/hellodhlyn/undersky"
)

var migration201904230001 = &gormigrate.Migration{
	ID: "201904230001",
	Migrate: func(db *gorm.DB) error {
		db.CreateTable()
		return db.CreateTable(
			us.Credential{},
		).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.DropTable(
			us.Credential{},
		).Error
	},
}
