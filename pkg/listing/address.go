package listing

// Address defines the properties of a address
type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Phone     string `json:"phone"`
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	Line3     string `json:"line3"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Zip       string `json:"zip"`
}
