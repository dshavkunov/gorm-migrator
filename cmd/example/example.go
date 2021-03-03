package main

import (
	"github.com/sirupsen/logrus"
	migrator "github.com/sulimak-co/gorm-migrator"
	"github.com/sulimak-co/gorm-migrator/cmd/example/migrations"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(nil, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mgr := migrator.NewMigrator(db, logrus.New(), migrations.GetMigrationsList())
	if err := mgr.Execute(); err != nil {
		panic(err)
	}
}
