package handlers

import (
	"html/template"
	"library-management/database"
	"library-management/models"
	"net/http"
	"strconv"
)

func BooksPage(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	rows, err := database.DB.Query("SELECT id, title, author, isbn, quantity FROM books ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Kitaplar getirilemedi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var book models.Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Quantity)
		if err != nil {
			http.Error(w, "Kitap okuma hatası", http.StatusInternalServerError)
			return
		}

		books = append(books, book)
	}

	tmpl := template.Must(template.ParseFiles("templates/books.html"))
	tmpl.Execute(w, books)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/books", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	isbn := r.FormValue("isbn")
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	_, err := database.DB.Exec(
		"INSERT INTO books (title, author, isbn, quantity) VALUES (?, ?, ?, ?)",
		title,
		author,
		isbn,
		quantity,
	)

	if err != nil {
		http.Error(w, "Kitap eklenemedi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	id := r.URL.Query().Get("id")

	_, err := database.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Kitap silinemedi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
