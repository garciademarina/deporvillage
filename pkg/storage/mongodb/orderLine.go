package mongodb

// OrderLines defines the storage form of an order
type OrderLines struct {
	SKU      string `json:"sku"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
}
