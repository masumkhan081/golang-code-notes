// query_timeout_and_tx.go
// Demonstrates query timeouts and transactions with database/sql + PostgreSQL.
// Requires a running Postgres instance.
// Usage: export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
//        go run query_timeout_and_tx.go
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"
)

func openDB() (*sql.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return nil, errors.New("DATABASE_URL is required")
    }

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(30 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}

func getTaskTitleByID(ctx context.Context, db *sql.DB, id int64) (string, error) {
    var title string
    err := db.QueryRowContext(ctx, `SELECT title FROM tasks WHERE id = $1`, id).Scan(&title)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", fmt.Errorf("task %d: %w", id, sql.ErrNoRows)
        }
        return "", err
    }
    return title, nil
}

func assignTask(ctx context.Context, db *sql.DB, taskID, employeeID int64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        _ = tx.Rollback()
    }()

    if _, err := tx.ExecContext(
        ctx,
        `UPDATE tasks SET assigned_to_user_id = $1 WHERE id = $2`,
        employeeID,
        taskID,
    ); err != nil {
        return err
    }

    return tx.Commit()
}

func main() {
    db, err := openDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    title, err := getTaskTitleByID(ctx, db, 1)
    if err != nil {
        log.Println("getTaskTitleByID:", err)
    } else {
        log.Println("task title:", title)
    }

    if err := assignTask(ctx, db, 1, 42); err != nil {
        log.Println("assignTask:", err)
    }
}
