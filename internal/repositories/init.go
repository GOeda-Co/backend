package repositories

import (
	"log"
	"os"
	"time"

	"github.com/pressly/goose/v3"

	"repeatro/internal/models"
	"repeatro/migrations"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(config *viper.Viper) *gorm.DB {
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

	db, err := gorm.Open(postgres.Open(config.GetString("database.connection_string")), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error during opening database")
	}

	db.AutoMigrate(&models.Card{}, &models.User{})

	return db
}

func InitGooseMigration(db *gorm.DB) {
	goose.SetBaseFS(migrations.EmbedMigrations)

	// Automatically fetch migrations
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	trueDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if err := goose.Up(trueDB, "."); err != nil {
		panic(err)
	}

	// connectionString := config.GetString("database.connection_string")
	// max_idle_connections := config.GetString("database.max_idle_connections")
	// max_open_connections := config.GetString("database.max_open_connections")
	// connection_max_lifetime := config.GetString("database.conecction_max_lifetime")
	// driver_name := config.GetString("database.driver_name")
	// if connectionString == "" {
	// 	log.Fatalf("Database connection string is missing")
	// }
	// dbHandler, err := sql.Open(driver_name, connectionString)
	// if err != nil {
	// 	log.Fatalf("Error %v", err)
	// }
	// max_idle_connections_int, err := strconv.Atoi(max_idle_connections)
	// if err != nil {
	// 	log.Fatalf("Error during conversion")
	// }

	// max_open_connections_int, err := strconv.Atoi(max_open_connections)
	// if err != nil {
	// 	log.Fatalf("Error during conversion")
	// }

	// connection_max_lifetime_duration, err := str2duration.ParseDuration(connection_max_lifetime)
	// if err != nil {
	// 	log.Fatalf("Error during conversion")
	// }

	// dbHandler.SetMaxIdleConns(max_idle_connections_int)
	// dbHandler.SetMaxOpenConns(max_open_connections_int)
	// dbHandler.SetConnMaxLifetime(connection_max_lifetime_duration)

	// err = dbHandler.Ping()
	// if err != nil {
	// 	dbHandler.Close()
	// 	log.Fatalf("Error while validatin base %v", err)
	// }
	// return dbHandler
}
