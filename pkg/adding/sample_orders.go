package adding

var DefaultOrders = []Order{
	{
		ID:     1,
		Amount: 10,
		Status: "init",
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
	},
}
