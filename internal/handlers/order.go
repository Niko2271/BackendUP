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

func RegisterOrderRoutes(r chi.Router, db *sql.DB) {
	//маленькая r просто сокращение Route
	r.Route("/order", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateOrder(w, r, db) })
		r.Get("/", func(w http.ResponseWriter, r *http.Request) { GetOrders(w, r, db) })
		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) { UpdateOrderStatus(w, r, db) })
	})

	// ну это чтобы с товарами в заказе робить
	r.Route("/order_item", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateOrderItem(w, r, db) })
		r.Get("/{order_id}", func(w http.ResponseWriter, r *http.Request) { GetOrderItems(w, r, db) })
	})
}

func CreateOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var o models.Order
	err := json.NewDecoder(r.Body).Decode(&o) // эта строка крч переводит json на go
	if err != nil {
		http.Error(w, "Не верный json", http.StatusBadRequest)
		return
	}
	if o.Id == "" {
		o.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	_, err = db.Exec("Insert into orders(id, userId, datetime, status) values(?, ?, ?, ?)", o.Id, o.UserId, o.Datetime, o.Status)
	if err != nil {
		http.Error(w, "Ошибка добавления в корзину: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Заказ оформлен")
}

func GetOrders(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	limit := 10
	offset := 0
	rows, err := db.Query("select id, userId, datetime, status from orders limit ? offset ?", limit, offset)
	if err != nil {
		http.Error(w, "ошибка бдшки", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.Id, &o.UserId, &o.Datetime, &o.Status)
		orders = append(orders, o)
	}
	w.Header().Set("Content-Type", "application/json") //формат данных
	json.NewEncoder(w).Encode(orders)                  // это наоборот с go на json после отправляем клиенту
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var o models.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}
	validStatuses := map[string]bool{
		"Принят": true, "В работе": true, "Собран": true, "Готов к выдаче": true, "Отменён": true,
	}
	if !validStatuses[o.Status] {
		http.Error(w, "Неверный статус заказа. Допустимые: Принят, В работе, Собран, Готов к выдаче, Отменён", http.StatusBadRequest)
		return
	}
	res, err := db.Exec("update orders set status = ? where id = ?", o.Status, id)
	if err != nil {
		http.Error(w, "ошибка бд", http.StatusInternalServerError)
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Заказа нема", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Статус заказа обновлен")
}

func CreateOrderItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var oi models.OrderItem
	err := json.NewDecoder(r.Body).Decode(&oi)
	if err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}
	if oi.Id == "" {
		oi.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	_, err = db.Exec("insert into order_items(id, orderId, productId, price) values(?, ?, ?, ?)", oi.Id, oi.OrderId, oi.ProductId, oi.Price)
	if err != nil {
		http.Error(w, "Ошибка добавления элемента: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Элемент заказа добавлен")
}

func GetOrderItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	orderId := chi.URLParam(r, "order_id")
	rows, err := db.Query("select id, orderId, productId, price from order_items where orderId = ?", orderId)
	if err != nil {
		http.Error(w, "Ошибка базы", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var oi models.OrderItem
		rows.Scan(&oi.Id, &oi.OrderId, &oi.ProductId, &oi.Price)
		items = append(items, oi)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
