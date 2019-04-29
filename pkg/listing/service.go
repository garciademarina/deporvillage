package listing

import "errors"

// ErrNotFound is used when a order could not be found.
var ErrNotFound = errors.New("order not found")

// Repository provides access to order repository.
type Repository interface {
	// AddOrder saves a given order to the repository.
	GetOrder(int64) (Order, error)
}

// Service provides order adding operations.
type Service interface {
	// GetOrder returns the order with given ID
	GetOrder(int64) (Order, error)
}

type service struct {
	r Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// GetOrder adds the given order(s) to the database
func (s *service) GetOrder(id int64) (Order, error) {
	o, err := s.r.GetOrder(id)
	if err != nil {
		return Order{}, err
	}

	return o, nil
}
