package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/garciademarina/deporvillage/pkg/adding"
	"github.com/garciademarina/deporvillage/pkg/listing"
	"github.com/garciademarina/deporvillage/pkg/updating"
	"github.com/go-chi/chi"
)

// OrderHandler handler struct for Order endpoints
type OrderHandler struct {
	l listing.Service
	a adding.Service
	u updating.Service
}

// NewOrderHandler create a new OrderHandler
func NewOrderHandler(l listing.Service, a adding.Service, u updating.Service) OrderHandler {
	return OrderHandler{
		l,
		a,
		u,
	}
}

// GetOrder handles GET /order requests.
func (h *OrderHandler) GetOrder(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] OrderHandler.GetOrder\n")
		if ID := chi.URLParam(r, "ID"); ID != "" {
			ID, err := strconv.ParseInt(ID, 10, 64)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
				return
			}
			account, err := h.l.GetOrder(ID)
			if err != nil {
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
				return
			}

			fmt.Printf("Accccccounts %v\n", account)
			respondwithJSON(w, http.StatusOK, &account)
		}

	}
}

// UpdateStatusOrder handles PUT /order requests.
func (h *OrderHandler) UpdateStatusOrder(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] OrderHandler.UpdateStatusOrder\n")

		var newOrder updating.Order
		err := decodeUpdateOrder(r, &newOrder)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
			return
		}

		err = h.u.UpdateStatusOrder(newOrder)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
			return
		}

		type success struct {
			Success string `json:"success"`
		}
		respondwithJSON(w, http.StatusOK, success{"ok"})
	}
}

func decodeAddOrder(r *http.Request, order *adding.Order) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(order)
	if err != nil {
		return errors.New("Cannot decode Order from body")
	}

	if order.ID == 0 {
		return errors.New("Order ID field not found")
	}
	return nil
}

// AddOrder handles POST /order requests.
func (h *OrderHandler) AddOrder(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] OrderHandler.AddOrder\n")

		var newOrder adding.Order
		err := decodeAddOrder(r, &newOrder)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
			return
		}

		err = h.a.AddOrder(newOrder)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Message: err.Error()})
			return
		}

		type success struct {
			Success string `json:"success"`
		}
		respondwithJSON(w, http.StatusOK, success{"ok"})
	}
}

func decodeUpdateOrder(r *http.Request, order *updating.Order) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(order)
	if err != nil {
		return errors.New("Cannot decode Order from body")
	}

	if order.ID == nil {
		return errors.New("Order 'ID' field not found")
	}
	if order.Status == nil {
		return errors.New("Order 'status' field not found")
	}

	return nil
}
