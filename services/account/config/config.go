package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hdkef/be-assignment/pkg/config"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type AccountConfig struct {
	AppName          string
	AppPort          string
	APIDomain        string
	WebDomain        string
	SuperTokenUrl    string
	SuperTokenAPIKey string
}

func InitAccountConfig() *AccountConfig {
	return &AccountConfig{
		AppName:          os.Getenv("APP_NAME"),
		AppPort:          os.Getenv("APP_PORT"),
		APIDomain:        os.Getenv("API_DOMAIN"),
		WebDomain:        os.Getenv("WEB_DOMAIN"),
		SuperTokenUrl:    os.Getenv("SUPER_TOKEN_URL"),
		SuperTokenAPIKey: os.Getenv("SUPER_TOKEN_API_KEY"),
	}
}

func InitDB() *sql.DB {
	config := config.InitPostgreConfig()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.Schema)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	return db
}

func InitRBMQ() *amqp.Connection {
	config := config.InitRBMQConfig()

	// Establish connection to RabbitMQ
	conn, err := amqp.Dial(config.RBMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return conn
}
