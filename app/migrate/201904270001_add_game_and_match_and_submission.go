package main

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"

	us "github.com/hellodhlyn/undersky"
)

var migration201904270001 = &gormigrate.Migration{
	ID: "201904270001",
	Migrate: func(db *gorm.DB) error {
		return db.CreateTable(
			us.Game{},
			us.Match{},
			us.Submission{},
		).Error
	},
	Rollback: func(db *gorm.DB) error {
		return db.DropTable(
			us.Game{},
			us.Match{},
			us.Submission{},
		).Error
	},
}
