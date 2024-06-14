package handlers

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Create an order")
}

func (o *Order) List(writer http.ResponseWriter, response *http.Request) {
	fmt.Println("List all orders")
}

func (o *Order) GetByID(writer http.ResponseWriter, response *http.Request) {
	fmt.Println("Get an order by ID")
}

func (o *Order) UpdateByID(writer http.ResponseWriter, response *http.Request) {
	fmt.Println("Update an order by ID")
}

func (o *Order) DeleteByID(writer http.ResponseWriter, response *http.Request) {
	fmt.Println("Delete order by ID")
}
