package dbsetup

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var err error

func ConnectDB() *pgxpool.Pool {
	db_url := os.Getenv("DB_URL")

	DB, err = pgxpool.New(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("database connected successfully")

	return DB
}