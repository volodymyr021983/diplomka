package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb() *DbContainer {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error during .env file loading")
	}
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Connecting to db failed")
		log.Fatal()
	}
	if db != nil {
		fmt.Println(db.Migrator().CurrentDatabase())
		fmt.Println("CONNECTED TO ServersAndUsers DATABASE")
	}

	dbContainer := &DbContainer{
		DB: db,
	}
	if dbContainer.DB == nil {
		log.Fatal("DbContainer is nil")
	}

	return dbContainer
}

type DbContainer struct {
	DB *gorm.DB
}
