package adding

// Order defines the storage form of an order
type Order struct {
	ID              int64        `json:"id"`
	Amount          int64        `json:"amount"`
	Status          string       `json:"status"`
	OrderLines      []OrderLines `json:"order_lines"`
	ShippingAddress Address      `json:"shipping_address"`
	BillingAddress  Address      `json:"billing_address"`
}
