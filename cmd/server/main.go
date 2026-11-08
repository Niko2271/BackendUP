package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"clotheShop/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Файл .env не найден, используем стандартные настройки")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5050"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "store.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	db.SetMaxOpenConns(1)
	if err != nil {
		slog.Error("Ошибка подключения к БД", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(`
    create table if not exists roles(
        id text primary key,
        name text not null unique
    );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    insert or ignore into roles (id, name) values 
        ('1', 'admin'),
        ('2', 'manager'),
        ('3', 'client');
    `)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists users(
        id text primary key,
        name text not null,
        email text not null unique,
        phone text,
        roleId text not null,
        registrationDate datetime default current_timestamp,
        foreign key(roleId) references roles(id)
    );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists products(
        id text primary key,
        name text not null,
        type text,
        gender text,
        size text,
        color text,
        description text,
        brand text,
        season text,
        price real not null 
    );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists carts(
        id text primary key,
        userId text not null,
        productId text not null,
        count integer not null default 1,
        foreign key(userId) references users(id),
        foreign key(productId) references products(id)
    )`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists orders(
        id text primary key,
        userId text not null,
        datetime text not null,
        status text not null,
        foreign key(userId) references users(id)
    )`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists order_items(
        id text primary key,
        orderId text not null,
        productId text not null,
        price real not null,
        foreign key(orderId) references orders(id),
        foreign key(productId) references products(id)
    );`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
    create table if not exists photos(
        id text primary key,
        filename text not null,
        productId text not null,
        foreign key(productId) references products(id)
    )`)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	handlers.RegisterProductRoutes(r, db)
	handlers.RegisterUserRoutes(r, db)
	handlers.RegisterCartRoutes(r, db)
	handlers.RegisterOrderRoutes(r, db)
	handlers.RegisterPhotoRoutes(r, db)
	handlers.RegisterRoleRoutes(r, db)

	slog.Info("Сервер успешно запущен", "port", port, "db_path", dbPath)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		slog.Error("Ошибка запуска сервера", "error", err)
	}
}
