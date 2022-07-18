package main

import (
	_ "embed"
	"github.com/flyflyhe/dbMigrate/internal/db"
	"log"
)

//go:embed ../../trafficLog/dsn0.txt
var dsn0 string

//go:embed ../../trafficLog/dsn1.txt
var dsn1 string

func main() {
	err := db.Migrate(dsn0, dsn1, "user", "user_bak", "mysql", "mysql")
	log.Println(err)
}
