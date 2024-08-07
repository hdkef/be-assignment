package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	cfg2 "github.com/hdkef/be-assignment/pkg/config"
	"github.com/hdkef/be-assignment/services/migrations/config"
	_ "github.com/lib/pq"
)

func main() {

	cfgPg := cfg2.InitPostgreConfig()

	cfg := config.InitMigrationConfig()

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfgPg.User, cfgPg.Password, cfgPg.Host, cfgPg.Port, cfgPg.DBName))
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigFilePath,
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil {
		panic(err)
	}
}
