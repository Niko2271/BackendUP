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

func RegisterCartRoutes(r chi.Router, db *sql.DB) {
	r.Route("/cart", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateCart(w, r, db) })
		r.Get("/{userId}", func(w http.ResponseWriter, r *http.Request) { GetCart(w, r, db) })
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { DeleteCart(w, r, db) })
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { UpdateCart(w, r, db) })
	})
}

func CreateCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var c models.Cart
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "Не верный json", http.StatusBadRequest)
		return
	}
	if c.Id == "" {
		c.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if c.Count == 0 {
		c.Count = 1
	}
	_, err = db.Exec("Insert into carts(id, userId, productId, count) values(?, ?, ?, ?)", c.Id, c.UserId, c.ProductId, c.Count)
	if err != nil {
		http.Error(w, "Ошибка добавления в корзину: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Товар добавлен в корзину")
}

func GetCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId := chi.URLParam(r, "userId")
	rows, err := db.Query("select id, userId, productId, count from carts where userId = ?", userId)
	if err != nil {
		http.Error(w, "Ошибка бдшки", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var carts []models.Cart
	for rows.Next() {
		var c models.Cart
		rows.Scan(&c.Id, &c.UserId, &c.ProductId, &c.Count)
		carts = append(carts, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carts)
}

func DeleteCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	res, err := db.Exec("delete from carts where id = ?", id)
	if err != nil {
		http.Error(w, "ошибка бд", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, " товар с id  не найден", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Товар удален")
}

func UpdateCart(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var c models.Cart
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, "ОШибка в запросе json", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("update carts set count = ? where id = ?", c.Count, id)
	if err != nil {
		http.Error(w, "Ошибка обновления корзины, ошибка базы данных", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Корзина не найдена", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Данные в корзине обновлены")
}
