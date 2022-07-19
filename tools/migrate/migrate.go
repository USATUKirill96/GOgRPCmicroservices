package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Provide version (int value)")
	}
	var postgresURL string
	if len(os.Args) > 2 && os.Args[2] == "test" {
		postgresURL = fmt.Sprintf("%s?sslmode=disable", os.Getenv("POSTGRES_DB_TEST_URL"))
	} else {
		postgresURL = fmt.Sprintf("%s?sslmode=disable", os.Getenv("POSTGRES_DB_URL"))
	}
	m, err := migrate.New(fmt.Sprintf("file://users/migrations"), postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	if os.Args[1] == "up" {
		err := m.Up()
		if err != nil {
			log.Fatal(err)
		}
	} else if os.Args[1] == "drop" {
		err := m.Drop()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		version, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		err = m.Migrate(uint(version))
		if err != nil {
			log.Fatal(err)
		}
	}

}
