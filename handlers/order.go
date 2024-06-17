package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"project/model"
	"project/repository/order"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo
}

func (o *Order) Create(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}
	err := o.Repo.Insert(request.Context(), order)
	if err != nil {
		fmt.Println("Failed to insert: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(res)
	writer.WriteHeader(http.StatusCreated)

	fmt.Println("Create an order")
}

func (o *Order) List(writer http.ResponseWriter, request *http.Request) {
	cursorStr := request.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	res, err := o.Repo.FindAll(request.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("Failed to find all:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:items`
		Next  uint64        `json:next,omitempty`
	}
	response.Items = res.Orders
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(data)
}

func (o *Order) GetByID(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")

	const base = 10
	const bitSize = 64
	orderId, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	orderObject, err := o.Repo.FindByID(request.Context(), orderId)
	if errors.Is(err, order.ErrNotExist) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Failed to find order with id: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(orderObject); err != nil {
		fmt.Println("failed to marshal:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Get an order by ID")
}

func (o *Order) UpdateByID(writer http.ResponseWriter, request *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(request, "id")

	const base = 10
	const bitSize = 64
	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := o.Repo.FindByID(request.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Failed to find by id", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.Update(request.Context(), theOrder)
	if err != nil {
		fmt.Println("Failed to marshal", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("Update an order by ID")
}

func (o *Order) DeleteByID(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(request.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Failed to delete by id: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("Delete order by ID")
}
