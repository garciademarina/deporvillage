package updating

// Order defines the update form of an order
type Order struct {
	ID     *int64  `json:"id"`
	Status *string `json:"status"`
}
