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

func RegisterPhotoRoutes(r chi.Router, db *sql.DB) {
	r.Route("/photo", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreatePhoto(w, r, db) })
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { GetPhoto(w, r, db) })
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { DeletePhoto(w, r, db) })
	})
}

func CreatePhoto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var p models.Photo
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}
	if p.Id == "" {
		p.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	_, err = db.Exec("insert into photos(id, filename, productId) values(?, ?, ?)", p.Id, p.Filename, p.ProductId)
	if err != nil {
		http.Error(w, "Ошбка добавления фото "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Фото добавлено")
}

func GetPhoto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var p models.Photo
	err := db.QueryRow("select id, filename, productId from photos where id = ?", id).Scan(&p.Id, &p.Filename, &p.ProductId)
	if err != nil {
		http.Error(w, "фото не найдено", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func DeletePhoto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	res, err := db.Exec("delete from photos where id = ?", id)
	if err != nil {
		http.Error(w, "ошибка бд", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "фото не найдено", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Фото удалено")
}
