package inits

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var conn *pgxpool.Pool

func InitDB(userName, password, host string, port int, dbname string) {
	ctx := context.Background()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", userName, password, host, port, dbname)
	var err error
	conn, err = pgxpool.New(ctx, connStr)

	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *pgxpool.Pool {
	return conn
}
func CloseDB(ctx context.Context) {
	conn.Close()
}
