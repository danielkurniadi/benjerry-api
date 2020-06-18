package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "github.com/iqdf/benjerry-service/config"
	"github.com/iqdf/benjerry-service/domain"
	productHTTP "github.com/iqdf/benjerry-service/product/delivery/http"
	productMongo "github.com/iqdf/benjerry-service/product/repository/mongo"
	"github.com/iqdf/benjerry-service/product/service"
)

const version = "1.0.0"
const usage string = `Ben Jerry Service.
Usage:
	app run [--port=<port>] [--host=<host>]
	app -h | --help
	app --version
Options:
	-h --help          Show this screen.
	--port=<port>      Set port where instance run.
	--host=<host>      Set hostname where instance run.`

// Command ...
type Command struct {
	Run     bool
	Port    int    `docopt:"--port"`
	Host    string `docopt:"--host"`
	Version bool
}

// parseCommand ...
func parseCommand() Command {
	command := Command{}

	// retrieve args & options
	opts, _ := docopt.ParseDoc(usage)
	opts.Bind(&command)
	return command
}

func main() {
	var (
		err     error
		command Command
		// config        config.Config
		dbConn         *mongo.Client
		productRepo    domain.ProductRepository
		productService domain.ProductService
		rootRouter     *mux.Router
		productRouter  *mux.Router
	)

	command = parseCommand()

	if !command.Run {
		if command.Version {
			fmt.Printf("ben&jerry %s \n", version)
		}
		return
	}

	// Setup database connection here ...
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	dbConn, err = mongo.Connect(ctx, mongoOpt)

	if err != nil {
		panic("unable to connect to mongodb")
	}

	// Setup repositories here ...
	productRepo = productMongo.NewProductRepo(dbConn, "tutorialDB") // benjerry

	// Instantiate services here ...
	productService = service.NewProductService(productRepo)

	// Register routings here ...
	rootRouter = mux.NewRouter()
	productRouter = rootRouter.PathPrefix("/api/products").Subrouter()

	productHTTP.NewProductHandler(productService).Routes(productRouter)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      rootRouter,
	}

	fmt.Println("Starting server...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Handle shutdowns when quit via SIGINT (Ctrl+C)
	// Note: SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Shutting Down...")
	os.Exit(0)
}
