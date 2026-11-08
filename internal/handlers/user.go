package handlers

import (
	"clotheShop/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, db *sql.DB) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateUser(w, r, db) })
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { GetUsers(w, r, db) })
		r.Get("/role/{role_id}", func(w http.ResponseWriter, r *http.Request) { GetUsersByRole(w, r, db) })
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { UpdateUser(w, r, db) })
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Не верный запрос", http.StatusBadRequest)
		return
	}
	if u.Id == "" {
		u.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	_, err = db.Exec("Insert into users(id, name, email, phone, roleId) values(?, ?, ?, ?, ?)", u.Id, u.Name, u.Email, u.Phone, u.RoleId)
	if err != nil {
		http.Error(w, "Ошибка добавления пользователя,проверьте email мб он зареган", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Пользователь создан еее")
}

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	limit := 10
	offset := 0
	rows, err := db.Query("select id, name, email, phone, roleId, registrationDate from users limit ? offset ?", limit, offset)
	if err != nil {
		http.Error(w, "ошибка бдшки", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.RoleId, &u.RegistrationDate)
		users = append(users, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// тут крч поиск по ролям
func GetUsersByRole(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roleID := chi.URLParam(r, "role_id")
	rows, err := db.Query("select id, name, email, phone, roleId, registrationDate from users where roleId = ?", roleID)
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Phone, &u.RoleId, &u.RegistrationDate)
		users = append(users, u)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Не верный запрос", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("update users set name = ?, email = ?, phone = ?, roleId = ? where id = ?", u.Name, u.Email, u.Phone, u.RoleId, id)
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Пользователя нема", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Данные пользователя обновлены")
}
