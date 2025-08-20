package database

import (
    "os"
    "path/filepath"
    "database/sql"
    "sync"
    "errors"
    _ "github.com/glebarez/go-sqlite"
    "time"
)

var (
    db   *sql.DB
    once sync.Once
)

type Stats struct {
    WPM      float64
    Accuracy float64
    AddedAt  time.Time
}

// ConnectToDatabase initializes the connection if it hasn't been initialized
func connectToDatabase() (*sql.DB, error) {
    var err error
    once.Do(func() {
        path, e := getDBPath()
        if e != nil {
            err = e
            return
        }
        db, err = sql.Open("sqlite", path)
        if err != nil {
            return
        }

        _, err = createTable(db)
    })

    return db, err
}

func getDBPath() (string, error) {
    configDir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }

    appDir := filepath.Join(configDir, "typtea")
    // Create directory if it doesn't exist
    if err := os.MkdirAll(appDir, 0700); err != nil {
        return "", err
    }

    return filepath.Join(appDir, "typtea.db"), nil
}

// CreateTable ensures the stats table exists
func createTable(db *sql.DB) (sql.Result, error) {
    sqlStmt := `CREATE TABLE IF NOT EXISTS stats (
        id INTEGER PRIMARY KEY,
        wpm REAL NOT NULL,
        accuracy REAL NOT NULL,
        added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
    return db.Exec(sqlStmt)
}

// Insert adds a new stats entry
func InsertStats(stats *Stats) (int64, error) {
    db, err := connectToDatabase()
    if err != nil {
        return 0, err
    }

    sqlStmt := `INSERT INTO stats (wpm, accuracy) VALUES (?, ?);`
    result, err := db.Exec(sqlStmt, stats.WPM, stats.Accuracy)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

type MaxStatsData struct {
    MaxWPM     float64
    MaxAccuracy float64
}

func MaxStats() (*MaxStatsData, error) {
    db, err := connectToDatabase()
    if err != nil {
        return nil, err
    }

    sqlStmt := `SELECT MAX(wpm), MAX(accuracy) FROM stats;`
    var maxWPM sql.NullFloat64
    var maxAccuracy sql.NullFloat64

    err = db.QueryRow(sqlStmt).Scan(&maxWPM, &maxAccuracy)
    if err != nil {
        return nil, err
    }

    if !maxWPM.Valid || !maxAccuracy.Valid {
        return nil, errors.New("no records found")
    }

    return &MaxStatsData{
        MaxWPM:     maxWPM.Float64,
        MaxAccuracy: maxAccuracy.Float64,
    }, nil
}
