package handlers

import (
	"html/template"
	"library-management/database"
	"library-management/models"
	"net/http"
)

func MembersPage(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	rows, err := database.DB.Query("SELECT id, name, email, phone FROM members ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Üyeler getirilemedi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []models.Member

	for rows.Next() {
		var member models.Member

		err := rows.Scan(&member.ID, &member.Name, &member.Email, &member.Phone)
		if err != nil {
			http.Error(w, "Üye okuma hatası", http.StatusInternalServerError)
			return
		}

		members = append(members, member)
	}

	tmpl := template.Must(template.ParseFiles("templates/members.html"))
	tmpl.Execute(w, members)
}

func AddMember(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/members", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	_, err := database.DB.Exec(
		"INSERT INTO members (name, email, phone) VALUES (?, ?, ?)",
		name,
		email,
		phone,
	)

	if err != nil {
		http.Error(w, "Üye eklenemedi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/members", http.StatusSeeOther)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	if !RequireLogin(w, r) {
		return
	}

	id := r.URL.Query().Get("id")

	_, err := database.DB.Exec("DELETE FROM members WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Üye silinemedi", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/members", http.StatusSeeOther)
}
