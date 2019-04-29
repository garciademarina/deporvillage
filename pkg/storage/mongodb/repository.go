package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/garciademarina/deporvillage/pkg/adding"
	"github.com/garciademarina/deporvillage/pkg/listing"
	"github.com/garciademarina/deporvillage/pkg/updating"
	"github.com/tkanos/gonfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// CollectionOrder identifier for the mongodb collection of orders
	CollectionOrder = "orders"
)

var (
	// ErrContextCreated ..
	ErrContextCreated = errors.New("Context cannot be created")
)

// Config represents mongodb configuration
type Config struct {
	Port     int    `json:"port"`
	URL      string `json:"url"`
	Database string `json:"database"`
}

// NewConfig creates new config.
func NewConfig(
	port int,
	url string,
	database string,
) Config {
	return Config{
		Port:     port,
		URL:      url,
		Database: database,
	}
}

// NewConfigFromFile creates new config from json file
func NewConfigFromFile(configFile string) (*Config, error) {
	c := Config{}
	err := gonfig.GetConf(configFile, &c)
	if err != nil {
		return &Config{}, err
	}
	return &c, nil
}

// Storage stores order data in mongodb
type Storage struct {
	Config     *Config
	Client     *mongo.Client
	Collection *mongo.Collection
}

// NewStorage returns a new mongodb storage
func NewStorage(config *Config) (*Storage, error) {
	var err error

	s := Storage{
		Config: config,
	}

	uri := fmt.Sprintf("mongodb://%s:%d", s.Config.URL, s.Config.Port)
	s.Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = s.Client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("mongodb: CONNECTED\n")

	s.Collection = s.Client.Database(s.Config.Database).Collection(CollectionOrder)
	return &s, nil
}

// AddOrder saves the given order to the repository
func (s *Storage) AddOrder(o adding.Order) error {

	// find if Order already exist
	var exist Order
	filter := bson.M{"id": o.ID}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := s.Collection.FindOne(ctx, filter)
	err := result.Decode(&exist)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if exist.ID != 0 {

		return adding.ErrDuplicate
	}

	// transform from adding.Order to mongodb.Order
	newOrderLines := make([]OrderLines, 0, len(o.OrderLines))
	for _, l := range o.OrderLines {
		newOrderLines = append(newOrderLines, OrderLines{
			SKU:      l.SKU,
			Price:    l.Price,
			Quantity: l.Quantity,
		})
	}
	getAddress := func(o adding.Address) Address {
		return Address{
			FirstName: o.FirstName,
			LastName:  o.LastName,
			Email:     o.Email,
			Company:   o.Company,
			Phone:     o.Phone,
			Line1:     o.Line1,
			Line2:     o.Line2,
			Line3:     o.Line3,
			City:      o.City,
			Country:   o.Country,
			Zip:       o.Zip,
		}
	}
	newShippingAddress := getAddress(o.ShippingAddress)
	newBillingAddress := getAddress(o.BillingAddress)

	newOrder := Order{
		ID:              o.ID,
		Amount:          o.Amount,
		Status:          o.Status,
		OrderLines:      newOrderLines,
		ShippingAddress: newShippingAddress,
		BillingAddress:  newBillingAddress,
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := s.Collection.InsertOne(ctx, newOrder)
	if err != nil {
		return err
	}
	id := res.InsertedID
	fmt.Printf("inserted order %v\n", id)
	return nil
}

// GetOrder returns an order with the specified ID
func (s *Storage) GetOrder(id int64) (listing.Order, error) {
	var o Order

	filter := bson.M{"id": id}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := s.Collection.FindOne(ctx, filter).Decode(&o)
	if err != nil {
		return listing.Order{}, listing.ErrNotFound
	}

	newOrderLines := make([]listing.OrderLines, 0, len(o.OrderLines))
	for _, l := range o.OrderLines {
		newOrderLines = append(newOrderLines, listing.OrderLines{
			SKU:      l.SKU,
			Price:    l.Price,
			Quantity: l.Quantity,
		})
	}
	getAddress := func(o Address) listing.Address {
		return listing.Address{
			FirstName: o.FirstName,
			LastName:  o.LastName,
			Email:     o.Email,
			Company:   o.Company,
			Phone:     o.Phone,
			Line1:     o.Line1,
			Line2:     o.Line2,
			Line3:     o.Line3,
			City:      o.City,
			Country:   o.Country,
			Zip:       o.Zip,
		}
	}
	newShippingAddress := getAddress(o.ShippingAddress)
	newBillingAddress := getAddress(o.BillingAddress)

	order := listing.Order{
		ID:              o.ID,
		Amount:          o.Amount,
		Status:          o.Status,
		OrderLines:      newOrderLines,
		ShippingAddress: newShippingAddress,
		BillingAddress:  newBillingAddress,
	}
	return order, nil

}

// UpdateStatusOrder ...
func (s *Storage) UpdateStatusOrder(order updating.Order) error {
	// find if Order already exist
	var exist Order
	filter := bson.M{"id": order.ID}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := s.Collection.FindOne(ctx, filter)
	err := result.Decode(&exist)
	if err == mongo.ErrNoDocuments {
		return updating.ErrNotFound
	}
	if err != nil {
		return err
	}

	exist.Status = *order.Status

	// set filters and updates
	filterUpdate := bson.M{"id": order.ID}
	update := bson.M{"$set": bson.M{"status": order.Status}}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, errUpdate := s.Collection.UpdateOne(ctx, filterUpdate, update)
	if errUpdate != nil {
		log.Fatal(err)
	}
	return nil
}
