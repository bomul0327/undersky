package main

import (
	"fmt"
	"os"

	gormigrate "gopkg.in/gormigrate.v1"

	us "github.com/hellodhlyn/undersky"
)

func main() {
	m := gormigrate.New(us.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migration201904220001,
		migration201904230001,
	})

	switch os.Args[1] {
	case "up":
		err := m.Migrate()
		if err != nil {
			panic(err)
		}
		fmt.Println("Migration succeed!")

	case "down":
		err := m.RollbackLast()
		if err != nil {
			panic(err)
		}
		fmt.Println("Rollback succeed!")
	}
}
