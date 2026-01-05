package main

import (
	"backend/internal/data"
	"backend/internal/driver"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	DB     *sql.DB
	Models data.Models
}

func main() {
	var cfg Config
	cfg.Port = "8080"

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...", err)
	}
	defer db.SQL.Close()

	// Auto-migration for existing tables
    migrateDB(db.SQL)

	app := &Application{
		Config: cfg,
		DB:     db.SQL,
		Models: data.New(db.SQL),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: app.routes(),
	}

	log.Printf("Starting server on port %s", cfg.Port)
    err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func migrateDB(db *sql.DB) {
    // Create tables if they don't exist
    queries := []string{
        `CREATE TABLE IF NOT EXISTS Offices (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS Departments (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );`,
        `CREATE TABLE IF NOT EXISTS Associates (
            id INT AUTO_INCREMENT PRIMARY KEY,
            first_name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            title VARCHAR(255),
            department VARCHAR(255),
            office VARCHAR(255),
            status VARCHAR(50),
            start_date DATETIME,
            empl_status VARCHAR(50),
            salary INT,
            dob DATETIME,
            profile_picture VARCHAR(255),
            email VARCHAR(255) UNIQUE,
            password VARCHAR(255)
        );`,
        `CREATE TABLE IF NOT EXISTS Tasks (
            id INT AUTO_INCREMENT PRIMARY KEY,
            requester_id INT NOT NULL,
            task_name VARCHAR(255) NOT NULL,
            task_value VARCHAR(255) NOT NULL,
            reason TEXT,
            status VARCHAR(50) DEFAULT 'pending',
            target_value INT,
            approvers JSON,
            timestamp INT,
            comments TEXT
        );`,
        `CREATE TABLE IF NOT EXISTS Thanks (
            id INT AUTO_INCREMENT PRIMARY KEY,
            from_id INT,
            to_id INT,
            message TEXT,
            category VARCHAR(255),
            timestamp BIGINT,
            FOREIGN KEY (from_id) REFERENCES Associates(id),
            FOREIGN KEY (to_id) REFERENCES Associates(id)
        );`,
        `CREATE TABLE IF NOT EXISTS time_off_requests (
            id INT AUTO_INCREMENT PRIMARY KEY,
            associate_id INT NOT NULL,
            start_date DATETIME NOT NULL,
            end_date DATETIME NOT NULL,
            reason TEXT,
            status VARCHAR(50) DEFAULT 'Pending',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (associate_id) REFERENCES Associates(id)
        );`,
        `CREATE TABLE IF NOT EXISTS time_entries (
            id INT AUTO_INCREMENT PRIMARY KEY,
            associate_id INT NOT NULL,
            date DATETIME NOT NULL,
            hours DECIMAL(5, 2) NOT NULL,
            overtime_hours DECIMAL(5, 2) NOT NULL DEFAULT 0,
            comments TEXT,
            status VARCHAR(50) DEFAULT 'Approved',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (associate_id) REFERENCES Associates(id)
        );`,
        // Attempt to add status column if it doesn't exist (for existing tables)
        `ALTER TABLE time_entries ADD COLUMN status VARCHAR(50) DEFAULT 'Approved';`,
        `CREATE TABLE IF NOT EXISTS menu_permissions (
            id INT AUTO_INCREMENT PRIMARY KEY,
            menu_item VARCHAR(255) NOT NULL,
            permission_type VARCHAR(50) NOT NULL,
            permission_value VARCHAR(255),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
        `CREATE TABLE IF NOT EXISTS thanks_categories (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        );`,
    }

    for _, query := range queries {
        _, err := db.Exec(query)
        if err != nil {
            log.Printf("Error executing migration query: %v\nQuery: %s", err, query)
        }
    }
    
    // Seed basic data if empty
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM Offices").Scan(&count)
    if err == nil && count == 0 {
         db.Exec("INSERT INTO Offices (name) VALUES ('London'), ('New York'), ('Paris')")
    }
    err = db.QueryRow("SELECT COUNT(*) FROM Departments").Scan(&count)
    if err == nil && count == 0 {
         db.Exec("INSERT INTO Departments (name) VALUES ('IT'), ('HR'), ('Design'), ('Sales')")
    }

    // Seed menu permissions for Time Entry if not exists
    err = db.QueryRow("SELECT COUNT(*) FROM menu_permissions WHERE menu_item = 'Time Entry'").Scan(&count)
    if err == nil && count == 0 {
         db.Exec("INSERT INTO menu_permissions (menu_item, permission_type, permission_value) VALUES ('Time Entry', 'everyone', NULL)")
    }

    // Seed thanks categories if not exists
    err = db.QueryRow("SELECT COUNT(*) FROM thanks_categories").Scan(&count)
    if err == nil && count == 0 {
         db.Exec("INSERT INTO thanks_categories (name) VALUES ('Team Player'), ('Superhero'), ('Thank You!'), ('Knowledge')")
    }
}
