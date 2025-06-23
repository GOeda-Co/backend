package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"
	_ "net/http/pprof"
	"repeatro/internal/config"
	"repeatro/internal/security"
	"repeatro/src/pkg"
	"repeatro/src/pkg/discovery/consul"
	userHttp "repeatro/src/user/internal/controller/http"
	"repeatro/src/user/internal/repository/postgresql"
	"repeatro/src/user/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)
const serviceName = "users"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	security := security.Security{ExpirationDelta: 600 * time.Minute}
	security.GetKyes()

	config := config.InitConfig("config")

	log.Println("Starting the user service")
	repo := postgresql.NewPostgresRepo(config, newLogger)
	service := services.CreateNewService(repo, &security)
	ctrl := userHttp.CreateNewUserController(service, &security)
	
	router := gin.Default()
	
	router.Handle(http.MethodPost, "/register", ctrl.Register)
	router.Handle(http.MethodPost, "/login", ctrl.Login)

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
