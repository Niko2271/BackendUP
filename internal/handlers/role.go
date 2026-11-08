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

func RegisterRoleRoutes(r chi.Router, db *sql.DB) {
	r.Route("/role", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateRole(w, r, db) })
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { GetRoles(w, r, db) })
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { UpdateRole(w, r, db) })
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { DeleteRole(w, r, db) })
	})
}

func CreateRole(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var role models.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}

	if role.Id == "" {
		role.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	_, err = db.Exec("insert into roles(id, name) values(?, ?)", role.Id, role.Name)
	if err != nil {
		http.Error(w, "Ошибка добавления роли", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Роль добавлена")
}

func GetRoles(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("select id, name from roles")
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		rows.Scan(&role.Id, &role.Name)
		roles = append(roles, role)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func UpdateRole(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var role models.Role

	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("update roles set name = ? where id = ?", role.Name, id)
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Роль не найдена", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Роль обновлена")
}

func DeleteRole(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")

	res, err := db.Exec("delete from roles where id = ?", id)
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Роль не найдена", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Роль удалена")
}
