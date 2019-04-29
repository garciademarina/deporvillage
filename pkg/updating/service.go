package updating

import "errors"

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

// Service provides order updating operations.
type Service interface {
	UpdateStatusOrder(Order) error
}

// Repository provides access to order repository.
type Repository interface {
	// UpdateStatusOrder update a given order to the repository.
	UpdateStatusOrder(Order) error
}

type service struct {
	oR Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// UpdateStatusOrder adds the given order to the database
func (s *service) UpdateStatusOrder(o Order) error {
	if !validStatus(*o.Status) {
		return ErrInvalidStatus
	}
	err := s.oR.UpdateStatusOrder(o)
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
