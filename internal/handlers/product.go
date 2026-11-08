package handlers

import (
	"clotheShop/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterProductRoutes(r chi.Router, db *sql.DB) {
	r.Route("/product", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateProduct(w, r, db) })
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { GetProducts(w, r, db) })
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) { DeleteProduct(w, r, db) })
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) { GetProductByID(w, r, db) })
	})
}

func CreateProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("Insert into products(id, name, type, gender, size, color, description, brand, season, price) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.Type, p.Gender, p.Size, p.Color, p.Description, p.Brand, p.Season, p.Price)
	if err != nil {
		http.Error(w, "Ошибка бдшки", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Товар добавлен")
}

func GetProducts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// чекаем параметры из ссылки если нет они пусты
	color := r.URL.Query().Get("color")
	brand := r.URL.Query().Get("brand")
	gender := r.URL.Query().Get("gender")
	size := r.URL.Query().Get("size")
	season := r.URL.Query().Get("season")
	productType := r.URL.Query().Get("type")

	// пагинация и сортировка
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}

	sort := r.URL.Query().Get("sort")
	sortOrder := "ASC"
	if sort == "desc" {
		sortOrder = "DESC"
	}

	query := "SELECT id, name, type, gender, size, color, description, brand, season, price FROM products WHERE 1=1"
	var args []interface{}

	if color != "" {
		query += " AND color = ?"
		args = append(args, color)
	}
	if brand != "" {
		query += " AND brand = ?"
		args = append(args, brand)
	}
	if gender != "" {
		query += " AND gender = ?"
		args = append(args, gender)
	}
	if size != "" {
		query += " AND size = ?"
		args = append(args, size)
	}
	if season != "" {
		query += " AND season = ?"
		args = append(args, season)
	}
	if productType != "" {
		query += " AND type = ?"
		args = append(args, productType)
	}

	query += " ORDER BY price " + sortOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Ошибка запроса к базе: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.Id, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price)
		products = append(products, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	res, err := db.Exec("delete from products where id = ?", id)
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

func GetProductByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var p models.Product
	err := db.QueryRow("select id, name, type, gender, size, color, description, brand, season, price from products where id = ?", id).Scan(&p.Id, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price)
	if err != nil {
		http.Error(w, "Товар не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
