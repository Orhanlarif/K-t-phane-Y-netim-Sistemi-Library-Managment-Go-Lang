package handlers

import (
	"html/template"
	"library-management/database"
	"library-management/models"
	"net/http"
	"time"
)

type DashboardData struct {
	BookCount       int
	MemberCount     int
	LoanCount       int
	OutOfStockCount int
	Today           string
	RecentLoans     []models.Loan
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	var bookCount int
	var memberCount int
	var loanCount int
	var outOfStockCount int

	database.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&bookCount)

	database.DB.QueryRow("SELECT COUNT(*) FROM members").Scan(&memberCount)

	database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM loans
		WHERE status = 'Ödünç Verildi'
	`).Scan(&loanCount)

	database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM books
		WHERE quantity = 0
	`).Scan(&outOfStockCount)

	rows, err := database.DB.Query(`
		SELECT
			books.title,
			members.name,
			loans.loan_date,
			IFNULL(loans.return_date, '-'),
			loans.status
		FROM loans
		INNER JOIN books ON books.id = loans.book_id
		INNER JOIN members ON members.id = loans.member_id
		ORDER BY loans.id DESC
		LIMIT 5
	`)

	if err != nil {
		http.Error(w, "Veriler getirilemedi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recentLoans []models.Loan

	for rows.Next() {
		var loan models.Loan

		rows.Scan(
			&loan.BookTitle,
			&loan.MemberName,
			&loan.LoanDate,
			&loan.ReturnDate,
			&loan.Status,
		)

		recentLoans = append(recentLoans, loan)
	}

	data := DashboardData{
		BookCount:       bookCount,
		MemberCount:     memberCount,
		LoanCount:       loanCount,
		OutOfStockCount: outOfStockCount,
		Today:           time.Now().Format("02.01.2006"),
		RecentLoans:     recentLoans,
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, data)
}
