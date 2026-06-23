package main

import (
	"fmt"
	"library-management/database"
	"library-management/handlers"
	"net/http"
)

func main() {
	database.InitDB()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.LoginPost(w, r)
			return
		}

		handlers.LoginPage(w, r)
	})

	http.HandleFunc("/dashboard", handlers.Dashboard)
	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/books", handlers.BooksPage)
	http.HandleFunc("/books/add", handlers.AddBook)
	http.HandleFunc("/books/delete", handlers.DeleteBook)

	http.HandleFunc("/members", handlers.MembersPage)
	http.HandleFunc("/members/add", handlers.AddMember)
	http.HandleFunc("/members/delete", handlers.DeleteMember)

	http.HandleFunc("/loans", handlers.LoansPage)
	http.HandleFunc("/loans/add", handlers.AddLoan)
	http.HandleFunc("/loans/return", handlers.ReturnLoan)

	fmt.Println("Sunucu çalışıyor: http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
