package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() *sql.DB {
	fmt.Println("Database connected.")
	db, err := sql.Open("mysql", "mydb:mydb@tcp(172.17.0.2:3306)/mydb")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
