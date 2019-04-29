package adding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type StorageStub struct{}

// AddOrder stub implementation
func (s *StorageStub) AddOrder(o Order) error {
	return nil
}

type StorageInvalidStatusStub struct{}

// AddOrder stub implementation
func (s *StorageInvalidStatusStub) AddOrder(o Order) error {
	return ErrDuplicate
}

type BrokerStub struct{}

var EventSent bool

// Publish stub implementation
func (s *BrokerStub) Publish(body string) error {
	EventSent = true
	return nil
}

// TestAddOrder AddOrder test
func TestAddOrder(t *testing.T) {
	EventSent = false
	r := &StorageStub{}
	b := &BrokerStub{}
	service := NewService(r, b)

	o := Order{
		ID:     1,
		Amount: 10,
		Status: "shipped",
		OrderLines: []OrderLines{
			OrderLines{
				SKU:      "TROP-UA-PLAS-09",
				Price:    10,
				Quantity: 1,
			},
			OrderLines{
				SKU:      "TROP-NP-PLAS-65",
				Price:    10,
				Quantity: 2,
			},
			OrderLines{
				SKU:      "TROP-LT-PLAS-89",
				Price:    5,
				Quantity: 10,
			},
		},
		ShippingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
		BillingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
	}

	err := service.AddOrder(o)
	assert.Equal(t, EventSent, true, "Event sent")
	assert.Nil(t, err, "Order cannot be created")
}

// TestAddOrderInvalidStatus AddOrder test
func TestAddOrderInvalidStatus(t *testing.T) {
	EventSent = false
	r := &StorageStub{}
	b := &BrokerStub{}
	service := NewService(r, b)

	o := Order{
		ID:     1,
		Amount: 10,
		Status: "invalidStatus",
		OrderLines: []OrderLines{
			OrderLines{
				SKU:      "TROP-UA-PLAS-09",
				Price:    10,
				Quantity: 1,
			},
			OrderLines{
				SKU:      "TROP-NP-PLAS-65",
				Price:    10,
				Quantity: 2,
			},
			OrderLines{
				SKU:      "TROP-LT-PLAS-89",
				Price:    5,
				Quantity: 10,
			},
		},
		ShippingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
		BillingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
	}

	err := service.AddOrder(o)
	assert.Equal(t, EventSent, false, "Event sent")
	assert.Equal(t, err, ErrInvalidStatus, "Order with invalid status")
}

// TestAddOrderInvalidStatus AddOrder test
func TestAddOrderDuplicated(t *testing.T) {

	r := &StorageInvalidStatusStub{}
	b := &BrokerStub{}
	service := NewService(r, b)

	o := Order{
		ID:     1,
		Amount: 10,
		Status: "shipped",
		OrderLines: []OrderLines{
			OrderLines{
				SKU:      "TROP-UA-PLAS-09",
				Price:    10,
				Quantity: 1,
			},
			OrderLines{
				SKU:      "TROP-NP-PLAS-65",
				Price:    10,
				Quantity: 2,
			},
			OrderLines{
				SKU:      "TROP-LT-PLAS-89",
				Price:    5,
				Quantity: 10,
			},
		},
		ShippingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
		BillingAddress: Address{
			FirstName: "Jhon",
			LastName:  "Snow",
			Email:     "j.snow@example.com",
			Company:   "Acme",
			Phone:     "555000111",
			Line1:     "711-2880 Nulla St.",
			Line2:     "",
			Line3:     "",
			City:      "Mankato",
			Country:   "Mississippi",
			Zip:       "96522",
		},
	}

	err := service.AddOrder(o)
	assert.Equal(t, err, ErrDuplicate, "Order duplicated")
}
