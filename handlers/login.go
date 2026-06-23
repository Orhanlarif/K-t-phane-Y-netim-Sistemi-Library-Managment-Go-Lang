package handlers

import (
	"html/template"
	"library-management/auth"
	"library-management/database"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))

	data := map[string]string{
		"Error": "",
	}

	tmpl.Execute(w, data)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var hashedPassword string

	err := database.DB.QueryRow(
		"SELECT password FROM users WHERE username = ?",
		username,
	).Scan(&hashedPassword)

	if err != nil || !auth.CheckPassword(password, hashedPassword) {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))

		data := map[string]string{
			"Error": "Kullanıcı adı veya şifre hatalı.",
		}

		tmpl.Execute(w, data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: username,
		Path:  "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RequireLogin(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("user")

	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}

	return true
}
