package database

import (
	"database/sql"
	"fmt"
	"library-management/auth"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite", "library.db")
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Veritabanına ulaşılamadı:", err)
	}

	createTables()
	createDefaultAdmin()

	fmt.Println("Veritabanı bağlantısı başarılı.")
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	bookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		isbn TEXT,
		quantity INTEGER NOT NULL
	);`

	memberTable := `
	CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT,
		phone TEXT
	);`

	loanTable := `
	CREATE TABLE IF NOT EXISTS loans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		book_id INTEGER NOT NULL,
		member_id INTEGER NOT NULL,
		loan_date TEXT NOT NULL,
		return_date TEXT,
		status TEXT NOT NULL,
		FOREIGN KEY(book_id) REFERENCES books(id),
		FOREIGN KEY(member_id) REFERENCES members(id)
	);`

	tables := []string{userTable, bookTable, memberTable, loanTable}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatal("Tablo oluşturma hatası:", err)
		}
	}
}

func createDefaultAdmin() {
	var count int

	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin").Scan(&count)
	if err != nil {
		log.Fatal("Admin kontrol hatası:", err)
	}

	if count == 0 {
		hashedPassword, err := auth.HashPassword("1234")
		if err != nil {
			log.Fatal("Şifre oluşturma hatası:", err)
		}

		_, err = DB.Exec(
			"INSERT INTO users (username, password) VALUES (?, ?)",
			"admin",
			hashedPassword,
		)

		if err != nil {
			log.Fatal("Admin oluşturma hatası:", err)
		}

		fmt.Println("Varsayılan admin kullanıcısı oluşturuldu.")
	}
}
