package dal

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var once sync.Once
var db *sql.DB

func Connect() (*sql.DB, error) {

	// dbname := "work"
	var err error
	once.Do(func() {
		// Connection parameters
connStr := "postgresql://root@localhost:26257/work?sslmode=disable"

		// connection_string := "user=2810k dbname=work sslmode=disable"
		// fmt.Println(connection_string)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			fmt.Println("Database Connection err", err)
			return
		}
		err = db.Ping()
		if err != nil {
			log.Fatal("Error pinging the database:", err)
		}
		fmt.Println("Successfully connected to CockroachDB!")

	})
	return db, err
}

func LogAndQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println(query)
	return db.Query(query, args...)
}

func GetDB() *sql.DB {
	return db
}
