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
	userHttp "repeatro/src/card/internal/controller/http"
	"repeatro/src/card/internal/repository/postgresql"
	"repeatro/src/card/internal/service"
	"repeatro/src/pkg/discovery/consul"

	"repeatro/src/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

const serviceName = "cards"

// TODO: technically each microservice should have separated main and current one should be divided into three
func main() {
	var port int
	flag.IntVar(&port, "port", 8084, "API handler port")
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
	service := services.CreateNewService(repo)
	ctrl := userHttp.CreateNewController(service)

	default_router := gin.Default()
	// default_router.Use(security.AuthMiddleware())

	router := default_router.Group("/cards")

	router.Handle(http.MethodPost, "", ctrl.AddCard)
	router.Handle(http.MethodGet, "", ctrl.ReadAllCardsToLearn)
	router.Handle(http.MethodPut, "/:id", ctrl.UpdateCard)
	router.Handle(http.MethodDelete, "/:id", ctrl.DeleteCard)
	router.Handle(http.MethodPost, "/answers", ctrl.AddAnswers)

	if err := default_router.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
