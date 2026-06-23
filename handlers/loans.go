package handlers

import (
	"html/template"
	"library-management/database"
	"library-management/models"
	"net/http"
	"strconv"
	"time"
)

type LoanPageData struct {
	Books   []models.Book
	Members []models.Member
	Loans   []models.Loan
	Error   string
	Success string
}

func LoansPage(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	bookRows, err := database.DB.Query("SELECT id, title, quantity FROM books ORDER BY title ASC")
	if err != nil {
		http.Error(w, "Kitaplar getirilemedi", http.StatusInternalServerError)
		return
	}
	defer bookRows.Close()

	var books []models.Book

	for bookRows.Next() {
		var book models.Book
		bookRows.Scan(&book.ID, &book.Title, &book.Quantity)
		books = append(books, book)
	}

	memberRows, err := database.DB.Query("SELECT id, name FROM members ORDER BY name ASC")
	if err != nil {
		http.Error(w, "Üyeler getirilemedi", http.StatusInternalServerError)
		return
	}
	defer memberRows.Close()

	var members []models.Member

	for memberRows.Next() {
		var member models.Member
		memberRows.Scan(&member.ID, &member.Name)
		members = append(members, member)
	}

	loanRows, err := database.DB.Query(`
		SELECT
			loans.id,
			books.title,
			members.name,
			loans.loan_date,
			IFNULL(loans.return_date, '-'),
			loans.status
		FROM loans
		INNER JOIN books ON books.id = loans.book_id
		INNER JOIN members ON members.id = loans.member_id
		ORDER BY loans.id DESC
	`)
	if err != nil {
		http.Error(w, "Ödünç kayıtları getirilemedi", http.StatusInternalServerError)
		return
	}
	defer loanRows.Close()

	var loans []models.Loan

	for loanRows.Next() {
		var loan models.Loan

		loanRows.Scan(
			&loan.ID,
			&loan.BookTitle,
			&loan.MemberName,
			&loan.LoanDate,
			&loan.ReturnDate,
			&loan.Status,
		)

		loans = append(loans, loan)
	}

	data := LoanPageData{
		Books:   books,
		Members: members,
		Loans:   loans,
		Error:   errorMessage,
		Success: successMessage,
	}

	tmpl := template.Must(template.ParseFiles("templates/loans.html"))
	tmpl.Execute(w, data)
}

func AddLoan(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/loans", http.StatusSeeOther)
		return
	}

	bookID := r.FormValue("book_id")
	memberID := r.FormValue("member_id")

	var quantity int

	err := database.DB.QueryRow(
		"SELECT quantity FROM books WHERE id = ?",
		bookID,
	).Scan(&quantity)

	if err != nil {
		http.Redirect(w, r, "/loans?error=Kitap bulunamadı", http.StatusSeeOther)
		return
	}

	if quantity <= 0 {
		http.Redirect(w, r, "/loans?error=Bu kitap stokta yok", http.StatusSeeOther)
		return
	}

	loanDate := time.Now().Format("2006-01-02")

	_, err = database.DB.Exec(`
		INSERT INTO loans (book_id, member_id, loan_date, status)
		VALUES (?, ?, ?, ?)
	`,
		bookID,
		memberID,
		loanDate,
		"Ödünç Verildi",
	)

	if err != nil {
		http.Redirect(w, r, "/loans?error=Ödünç işlemi başarısız", http.StatusSeeOther)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE books SET quantity = quantity - 1 WHERE id = ?",
		bookID,
	)

	if err != nil {
		http.Redirect(w, r, "/loans?error=Stok güncellenemedi", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/loans?success=Kitap başarıyla ödünç verildi", http.StatusSeeOther)
}

func ReturnLoan(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	id := r.URL.Query().Get("id")

	loanID, err := strconv.Atoi(id)
	if err != nil {
		http.Redirect(w, r, "/loans?error=Geçersiz işlem", http.StatusSeeOther)
		return
	}

	var bookID int
	var status string

	err = database.DB.QueryRow(
		"SELECT book_id, status FROM loans WHERE id = ?",
		loanID,
	).Scan(&bookID, &status)

	if err != nil {
		http.Redirect(w, r, "/loans?error=Ödünç kaydı bulunamadı", http.StatusSeeOther)
		return
	}

	if status == "İade Edildi" {
		http.Redirect(w, r, "/loans?error=Bu kitap zaten iade edilmiş", http.StatusSeeOther)
		return
	}

	returnDate := time.Now().Format("2006-01-02")

	_, err = database.DB.Exec(`
		UPDATE loans
		SET status = ?, return_date = ?
		WHERE id = ?
	`,
		"İade Edildi",
		returnDate,
		loanID,
	)

	if err != nil {
		http.Redirect(w, r, "/loans?error=İade işlemi başarısız", http.StatusSeeOther)
		return
	}

	_, err = database.DB.Exec(
		"UPDATE books SET quantity = quantity + 1 WHERE id = ?",
		bookID,
	)

	if err != nil {
		http.Redirect(w, r, "/loans?error=Stok iadesi başarısız", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/loans?success=Kitap başarıyla iade alındı", http.StatusSeeOther)
}
