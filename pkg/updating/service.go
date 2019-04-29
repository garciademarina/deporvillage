package updating

import (
	"errors"
	"fmt"
)

// Payload ...
type Payload []Order

// Event defines possible outcomes from the "adding actor"
type Event int

const (
	// Done means finished processing successfully
	Done Event = iota

	// OrderNotFound means the given order has not been found
	OrderNotFound

	// Failed means processing did not finish successfully
	Failed
)

// ErrNotFound ...
var ErrNotFound = errors.New("Order not found")

// ErrInvalidStatus ...
var ErrInvalidStatus = errors.New("Order status is invalid")

// GetMeaning ...
func (e Event) GetMeaning() string {
	if e == Done {
		return "Done"
	}

	if e == OrderNotFound {
		return "Order not order"
	}

	if e == Failed {
		return "Failed"
	}

	return "Unknown result"
}

// EventOrderStatusUpdated ...
var EventOrderStatusUpdated = "order_status_updated"

// Service provides order updating operations.
type Service interface {
	UpdateStatusOrder(Order) error
}

// Repository provides access to order repository.
type Repository interface {
	// UpdateStatusOrder update a given order to the repository.
	UpdateStatusOrder(Order) error
}

// Broker provides access to a message broker
type Broker interface {
	// Publish sends a message
	Publish(body string) error
}

type service struct {
	oR Repository
	b  Broker
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository, b Broker) Service {
	return &service{r, b}
}

// UpdateStatusOrder adds the given order to the database
func (s *service) UpdateStatusOrder(o Order) error {
	if !validStatus(*o.Status) {
		return ErrInvalidStatus
	}
	err := s.oR.UpdateStatusOrder(o)

	// send event
	s.b.Publish(fmt.Sprintf("%d-%s [%s]", o.ID, EventOrderStatusUpdated, *o.Status))

	if err != nil {
		return err
	}
	return nil
}

func validStatus(s string) bool {
	set := make(map[string]bool)
	set["pending_confirmation"] = true
	set["confirmed"] = true
	set["sent_to_warehouse"] = true
	set["shipped"] = true
	set["in_transit"] = true
	set["delivered"] = true

	return set[s]
}
