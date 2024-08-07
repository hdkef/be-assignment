package config

import "os"

type MigrationConfig struct {
	MigFilePath string
}

func InitMigrationConfig() *MigrationConfig {

	return &MigrationConfig{
		MigFilePath: os.Getenv("MIGRATION_FILE_PATH"),
	}

}
