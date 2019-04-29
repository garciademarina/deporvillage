package adding

import "errors"

type Payload []Order

// Event defines possible outcomes from the "adding actor"
type Event int

const (
	// Done means finished processing successfully
	Done Event = iota

	// OrderAlreadyExists means the given order is a duplicate of an existing one
	OrderAlreadyExists

	// Failed means processing did not finish successfully
	Failed
)

// GetMeaning ...
func (e Event) GetMeaning() string {
	if e == Done {
		return "Done"
	}

	if e == OrderAlreadyExists {
		return "Duplicate order"
	}

	if e == Failed {
		return "Failed"
	}

	return "Unknown result"
}

// ErrDuplicate ...
var ErrDuplicate = errors.New("order already exists")

// ErrInvalidStatus ...
var ErrInvalidStatus = errors.New("Order status is invalid")

// Service provides order adding operations.
type Service interface {
	AddOrder(Order) error
	AddSampleOrders(Payload) <-chan Event
}

// Repository provides access to order repository.
type Repository interface {
	// AddOrder saves a given order to the repository.
	AddOrder(Order) error
}

type service struct {
	oR Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// AddOrder adds the given order to the database
func (s *service) AddOrder(o Order) error {
	if !validStatus(o.Status) {
		return ErrInvalidStatus
	}

	err := s.oR.AddOrder(o)
	if err != nil {
		return err
	}
	return nil
}

// AddSampleOrders adds some sample order to the database
func (s *service) AddSampleOrders(data Payload) <-chan Event {
	results := make(chan Event)

	go func() {
		defer close(results)

		for _, o := range data {
			err := s.oR.AddOrder(o)
			if err != nil {
				if err == ErrDuplicate {
					results <- OrderAlreadyExists
					continue
				}
				results <- Failed
				continue
			}

			results <- Done
		}
	}()

	return results
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
