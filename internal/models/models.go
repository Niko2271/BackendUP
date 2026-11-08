package models

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Gender      string  `json:"gender"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Season      string  `json:"season"`
	Price       float64 `json:"price"`
}

type User struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	RoleId           string `json:"roleId"`
	RegistrationDate string `json:"registrationDate"`
}

type Cart struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	ProductId string `json:"productId"`
	Count     int    `json:"count"`
}

type Order struct {
	Id       string `json:"id"`
	UserId   string `json:"userId"`
	Datetime string `json:"datetime"`
	Status   string `json:"status"`
}

type OrderItem struct {
	Id        string  `json:"id"`
	OrderId   string  `json:"orderId"`
	ProductId string  `json:"productId"`
	Price     float64 `json:"price"`
}

type Photo struct {
	Id        string `json:"id"`
	Filename  string `json:"filename"`
	ProductId string `json:"productId"`
}

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
