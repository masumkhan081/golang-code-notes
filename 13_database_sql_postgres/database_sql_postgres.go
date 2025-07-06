// database_sql_postgres.go
// Demonstrates connecting to PostgreSQL using database/sql.
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	// This example assumes you have PostgreSQL running and accessible.
	connStr := "user=postgres password=secret dbname=mydb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer db.Close()

	// Simple query
	var now string
	err = db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		fmt.Println("Query error:", err)
		return
	}
	fmt.Println("Current time from DB:", now)
}
// Documentation:
// - Use database/sql for DB access; pq for PostgreSQL driver.
// - Always check errors and close DB.
// - Edge cases: connection failure, query errors, credentials.
