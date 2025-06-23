package main

import (
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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

// TODO: technically each microservice should have separated main and current one should be divided into three
func main() {
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

	if err := default_router.Run(":8084"); err != nil {
		panic(err)
	}
}
