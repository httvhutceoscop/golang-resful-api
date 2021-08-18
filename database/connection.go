package database

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
)

func ConnectDB() (c *pgx.Conn, err error) {
	dbName := "go-api"
	dbUser := "postgres"
	dbPass := ""
	dbHost := "localhost"
	dbPort := "5432"
	conn, err := pgx.Connect(context.Background(), "postgresql://"+dbUser+":"+dbPass+"@"+dbHost+":"+dbPort+"/"+dbName)
	if err != nil {
		fmt.Println("Error connecting to DB")
		fmt.Println(err.Error())
	}
	_ = conn.Ping(context.Background())
	return conn, err
}
