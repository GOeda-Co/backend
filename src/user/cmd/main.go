package main

import (
	"log"
	"os"
	"time"

	_ "net/http/pprof"
	"repeatro/internal/config"
	"repeatro/internal/security"
	"net/http"
	userHttp "repeatro/src/user/internal/controller/http"
	"repeatro/src/user/internal/repository/postgresql"
	"repeatro/src/user/internal/service"

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
	service := services.CreateNewService(repo, &security)
	ctrl := userHttp.CreateNewUserController(service, &security)
	
	router := gin.Default()
	
	router.Handle(http.MethodPost, "/register", ctrl.Register)
	router.Handle(http.MethodPost, "/login", ctrl.Login)

	if err := router.Run(":8082"); err != nil {
		panic(err)
	}
}
