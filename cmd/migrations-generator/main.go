package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type migrationFile struct {
	name   string
	number int
}

var (
	migrationsPath = "migrations"
	listFile       = "migrations.go"
)

func main() {
	migrationName := flag.String("name", "", "migration name in lowerCamelCase")
	flag.Parse()

	if *migrationName == "" {
		panic(fmt.Errorf("migration name is emptry"))
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		if err := os.Mkdir(migrationsPath, 0775); err != nil {
			panic(err)
		}
	}

	createdMigrations, err := getCreatedMigrations()
	if err != nil {
		panic(err)
	}

	if err := checkMigrationExist(*migrationName, createdMigrations); err != nil {
		panic(err)
	}

	number := getMigrationNumber(createdMigrations)

	if err := createMigrationFile(*migrationName, number); err != nil {
		panic(err)
	}

	createdMigrations = append(createdMigrations, migrationFile{
		name:   *migrationName,
		number: number,
	})

	if err := updateAllMigrations(createdMigrations); err != nil {
		panic(err)
	}

	fmt.Println("migration successfully created")
}

func getCreatedMigrations() ([]migrationFile, error) {
	var allMigrationNames []migrationFile

	if err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
		fileName := info.Name()
		if fileName == migrationsPath || fileName == listFile {
			return nil
		}
		temp := regexp.MustCompile(`[_.]`)
		parts := temp.Split(fileName, -1)
		if len(parts) != 3 {
			return fmt.Errorf("wrong file name format: %v", fileName)
		}
		number, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		allMigrationNames = append(allMigrationNames, migrationFile{
			name:   parts[1],
			number: number,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return allMigrationNames, nil
}

func checkMigrationExist(name string, createdMigrations []migrationFile) error {
	for _, v := range createdMigrations {
		if name == v.name {
			return fmt.Errorf("migration with given name alreade exists")
		}
	}
	return nil
}

func getMigrationNumber(createdMigrations []migrationFile) int {
	if len(createdMigrations) == 0 {
		return 1
	}
	return createdMigrations[len(createdMigrations)-1].number + 1
}

func createMigrationFile(migrationName string, number int) error {
	t, err := template.New("migration").Parse(migrationTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%v/%v_%v.go", migrationsPath, fmt.Sprintf("%06d", number), migrationName))
	if err != nil {
		return err
	}

	data := struct {
		MigrationName string
	}{MigrationName: migrationName}

	return t.Execute(f, data)
}

func updateAllMigrations(createdMigrations []migrationFile) error {
	listTemp, err := template.New("list").Parse(migrationListTemplate)
	if err != nil {
		return err
	}

	itemTemp, err := template.New("item").Parse(migrationListItemTemplate)
	if err != nil {
		return err
	}

	var items string
	for i, v := range createdMigrations {
		var result bytes.Buffer
		data := struct {
			MigrationName string
		}{MigrationName: v.name}
		if err := itemTemp.Execute(&result, data); err != nil {
			return err
		}
		if i != 0 {
			items += "\n"
		}
		items += result.String()
	}

	f, err := os.Create(fmt.Sprintf("%v/%v", migrationsPath, listFile))
	if err != nil {
		return err
	}

	data := struct {
		Items string
	}{Items: items}

	return listTemp.Execute(f, data)
}
