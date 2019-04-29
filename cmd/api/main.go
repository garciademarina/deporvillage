package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/garciademarina/deporvillage/pkg/adding"
	broker "github.com/garciademarina/deporvillage/pkg/broker/rabbitmq"
	"github.com/garciademarina/deporvillage/pkg/listing"
	"github.com/garciademarina/deporvillage/pkg/server"
	"github.com/garciademarina/deporvillage/pkg/storage/mongodb"
	"github.com/garciademarina/deporvillage/pkg/updating"
)

var port = flag.Int("port", 8080, "-port=<port> sets the server's listening port. 8080 by default.")
var env = flag.String("env", "prod", "-env=<environment> specifies the environment. prod by default.")
var configFile = flag.String("conf", "config.json", "-conf=<configuration> specifies the configuration file path. config.json by default.")

func init() {
	flag.Parse()
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	configMongo, err := mongodb.NewConfigFromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	storage, err2 := mongodb.NewStorage(configMongo)
	if err2 != nil {
		log.Fatal(err2)
	}

	// broker rabbitmq
	broker, errBroker := broker.NewBroker("amqp://guest:guest@rabbitmq", "deporvillageQueue")
	if errBroker != nil {
		log.Fatal(errBroker)
	}
	defer broker.Close()

	adder := adding.NewService(storage, broker)
	updater := updating.NewService(storage, broker)
	lister := listing.NewService(storage)

	resultsOrders := adder.AddSampleOrders(adding.DefaultOrders)

	go func() {
		for result := range resultsOrders {
			fmt.Printf("Added sample beer with result %s.\n", result.GetMeaning()) // human-friendly
		}
	}()

	config := server.NewConfig(*port, *env)
	s := server.NewServer(config, logger, adder, lister, updater)

	err = s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ensureInterruptionsStopApplication(cancel, logger)
}

func ensureInterruptionsStopApplication(cancelFunc context.CancelFunc, logger *log.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Println(fmt.Sprintf("Got signal %s. Stopping server...", s))
		cancelFunc()

		os.Exit(1)
	}()
}
