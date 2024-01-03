package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// var DB *sqlx.DB

func Connect() *sqlx.DB {
	Db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_URI"))

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to DB")

	return Db
}
