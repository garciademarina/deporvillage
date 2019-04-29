package mongodb

// Order defines the storage form of an order
type Order struct {
	ID              int64        `json:"id" bson:"id"`
	Amount          int64        `json:"amount" bson:"amount"`
	Status          string       `json:"status" bson:"status"`
	OrderLines      []OrderLines `json:"order_lines" bson:"order_lines"`
	ShippingAddress Address      `json:"shipping_address" bson:"shipping_address"`
	BillingAddress  Address      `json:"billing_address" bson:"billing_address"`
}
