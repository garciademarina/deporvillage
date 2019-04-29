package mongodb

// Address defines the properties of a address
type Address struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Company   string `json:"company" bson:"company"`
	Phone     string `json:"phone" bson:"phone"`
	Line1     string `json:"line1" bson:"line1"`
	Line2     string `json:"line2" bson:"line2"`
	Line3     string `json:"line3" bson:"line3"`
	City      string `json:"city" bson:"city"`
	Country   string `json:"country" bson:"country"`
	Zip       string `json:"zip" bson:"zip"`
}
