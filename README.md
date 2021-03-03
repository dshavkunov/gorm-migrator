# Gorm-migrator

Custom migrations for gorm

### Usage

Build `migrations-generator`

Create migrations
```shell
migrations-generator -name=createUserTable
```

Implement created migrations using gorm migrator

```go
package migrations

import (
	"gorm.io/gorm"
)

func createUserTableUp(migrator gorm.Migrator) error {
	type user struct {
		ID string `json:"id" gorm:"type:uuid;primary_key"`
	}
	return migrator.CreateTable(&user{})
}

func createUserTableDown(migrator gorm.Migrator) error {
	type user struct{}
	return migrator.DropTable(&user{})
}
```

Execute migrations
```go
mgr := migrator.NewMigrator(dbConn, logger, migrations.GetMigrationsList())
if err := mgr.Execute(); err != nil {
	panic(err)
}
```
