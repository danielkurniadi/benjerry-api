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
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "github.com/iqdf/benjerry-service/config"
	"github.com/iqdf/benjerry-service/common/auth"
	"github.com/iqdf/benjerry-service/common/config"
	"github.com/iqdf/benjerry-service/common/middleware"
	"github.com/iqdf/benjerry-service/domain"

	productHTTP "github.com/iqdf/benjerry-service/product/delivery/http"
	productMongo "github.com/iqdf/benjerry-service/product/repository/mongo"

	userHTTP "github.com/iqdf/benjerry-service/user/delivery/http"
	userMongo "github.com/iqdf/benjerry-service/user/repository/mongo"

	productUC "github.com/iqdf/benjerry-service/product/service"
	userUC "github.com/iqdf/benjerry-service/user/service"
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
	Port    string `docopt:"--port"`
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
		dbConn      *mongo.Client
		productRepo domain.ProductRepository
		userRepo    domain.UserRepository

		productService domain.ProductService
		userService    domain.UserService
		authService    *auth.Service

		rootRouter    *mux.Router
		productRouter *mux.Router
		userRouter    *mux.Router
	)

	command = parseCommand()

	if !command.Run {
		if command.Version {
			fmt.Printf("ben&jerry %s \n", version)
		}
		return
	}

	appconfig := config.Get(config.BENJERRY, command.Host, command.Port)
	config.PrintConfig(appconfig)

	appname := string(appconfig.AppName)

	// Setup database connection here ...
	ctx, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelMongo()

	mongoOpt := options.Client().ApplyURI(appconfig.DatabaseURI)
	dbConn, err = mongo.Connect(ctx, mongoOpt)

	if err != nil {
		panic("unable to connect to mongodb")
	}

	ctx, cancelRedis := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelRedis()

	redisConn, err := redis.DialURL(appconfig.RedisURI)
	if err != nil {
		panic(err)
	}

	// Setup repositories here ...
	productRepo = productMongo.NewProductRepo(dbConn, appconfig.DatabaseName) // benjerry
	userRepo = userMongo.NewUserRepo(dbConn, appconfig.DatabaseName)

	// Instantiate services here ...
	productService = productUC.NewProductService(productRepo)
	userService = userUC.NewUserService(appname, userRepo)
	authService = auth.NewAuthService(redisConn)

	// Setup Middleware here ....
	authMiddleware := middleware.AuthMiddleWare(authService)
	roleMiddleware := middleware.RoleMiddleWare(appname)
	middlewareChain := alice.New(authMiddleware, roleMiddleware)

	// Register routings here ...
	rootRouter = mux.NewRouter()
	productRouter = rootRouter.PathPrefix("/api/products").Subrouter()
	userRouter = rootRouter.PathPrefix("/api/users").Subrouter()

	sessionExpiry := 480 * time.Second
	productHTTP.NewProductHandler(productService).Routes(productRouter, middlewareChain)
	userHTTP.NewUserHandler(userService, authService, sessionExpiry).Routes(userRouter)

	server := &http.Server{
		Addr:         appconfig.AppAddress(),
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

	ctx, cancelRun := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelRun()

	server.Shutdown(ctx)

	log.Println("Shutting Down...")
	os.Exit(0)
}
