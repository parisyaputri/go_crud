package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDb() (*sql.DB, error) {

	// Gunakan nilai variabel lingkungan untuk koneksi database
	dsn := "root:parisol@tcp(127.0.0.1:3306)/final?charset=utf8mb4&parseTime=True&loc=Local"
	// change this to your db

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to db:", err)
		return nil, err
	}

	// Cek koneksi database
	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping Error:", err)
		return nil, err
	}

	// Eksekusi perintah SQL untuk membuat tabel-tabel jika belum ada
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			ID INT AUTO_INCREMENT PRIMARY KEY,
			Email VARCHAR(255) UNIQUE NOT NULL,
			Role INT NOT NULL,
			Password VARCHAR(255) NOT NULL
		)
	`)

	if err != nil {
		log.Fatal("Error creating user table:", err)
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			ID INT AUTO_INCREMENT PRIMARY KEY,
			ProductName VARCHAR(255) NOT NULL,
			Price INT NOT NULL,
			Quantity INT NOT NULL,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating products table:", err)
		return nil, err
	}

	return db, nil
}
