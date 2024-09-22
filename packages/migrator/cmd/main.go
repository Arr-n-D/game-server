package main

import (
	"context"
	"flag"
	"log"
	"os"

	"migrator/database"
	_ "migrator/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	ctx := context.Background()
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	dir := os.Getenv("MIGRATIONS_FOLDER")
	table := "migrations"
	goose.SetTableName(table)

	db := database.ConnectToDatabase()

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.RunContext(ctx, command, db, dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
