package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"

	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/events"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/flows"
	"github.com/form3tech-oss/go-flow/internal/sample_app/payment-flow/api/settings"
	"github.com/giantswarm/retry-go"

	"github.com/form3tech/go-cdc/cdc"
	"github.com/form3tech/go-cqrs/v7/cqrs"
	"github.com/form3tech/go-logger/log"
	"github.com/jmoiron/sqlx"
	"github.com/liamg/waitforhttp"
	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/source/file"
	"github.com/spf13/viper"

	migrate "github.com/rubenv/sql-migrate"
)

func Configure() {
	viper.AutomaticEnv()
	viper.SetDefault("MessageVisibilityTimeout", 60)

	db := ConnectToDatabase()

	migrateDatabase(db.DB)

	setupRoutes(db)
}

func setupRoutes(db *sqlx.DB) {
	router := gin.New()
	router.Use(gin.Recovery())
		http.HandleFunc("/", router.ServeHTTP)


	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	v1 := router.Group("/v2")
	v1.POST("/payments", handlePayment(db))
}

func getDBConnectionString() connectionString {
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("database-username", "postgres")
	viper.SetDefault("database-password", "password")
	host := viper.GetString("DATABASE_HOST")
	username := getOrDefaultString("database-username", "payment")
	password := viper.GetString("database-password")
	sslMode := getOrDefaultString("database-ssl-mode", "disable")
	port := getOrDefaultInt("DATABASE_PORT", 5432)

	return newConnectionString(host, username, password, "postgres", port, sslMode)
}

func ConnectToDatabase() *sqlx.DB {
	var db *sqlx.DB
	var err error

	_ = retry.Do(func() error {
		db, err = sqlx.Connect("postgres", getDBConnectionString().String())
		if err != nil {
			return err
		}

		return nil

	}, retry.MaxTries(10), retry.Sleep(time.Duration(200*time.Millisecond)))

	if err != nil {
		panic(err)
	}

	return db
}

func migrateDatabase(db *sql.DB) {
	n, err := migrate.Exec(db, "postgres", &migrate.FileMigrationSource{Dir: "api/migrations"}, migrate.Up)
	if err != nil {
		panic(fmt.Sprintf("could not migrate database, error: %v", err))
	}

	log.Infof("applied %d database migrations!\n", n)
}

func StartServer(stopChannel <-chan os.Signal, startedSignal chan<- bool, stopped chan<- struct{}) error {
	connectionString := getDBConnectionString()
	changeDataCapture, err := cdc.NewChangeDataCapture(
		cdc.DatabaseConnection{
			Host:     connectionString.host,
			Port:     connectionString.port,
			User:     connectionString.user,
			Password: connectionString.password,
			Database: connectionString.database,
		},
		settings.CDCSlotName,
		flows.NewMessagingSource(),
		cqrs.SharedQueueNamingStrategy{ApplicationName: settings.ApiName, QueueName: "events"}.GetQueueNameFor,
		-1)

	if err != nil {
		return err
	}

	changeDataCapture.RegisterChangeEvent("Payment", "insert", events.PaymentCreatedEvent{})

	port := settings.ServerPort
	address := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: address, Handler: nil}

	go func() {
		if signal := <-stopChannel; signal != nil {
			log.Info("Shutting down as requested...")
			_ = server.Close()
		}
	}()

	go func() {
		err := waitforhttp.Wait(server, time.Second*10)
		startedSignal <- err == nil
		if err != nil {
			log.Errorf("Server failed to start: %s", err)
			return
		}
		log.Infof("Server started on %s", server.Addr)
	}()

	go func() {
		if err := changeDataCapture.StartListening(); err != nil {
			log.Fatalf("Error from cdc: %s", err.Error())
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		changeDataCapture.StopListening()
		stopped <- struct{}{}
		if err == http.ErrServerClosed {
			log.Info("Server was shut down cleanly.")
			return nil
		}
		return err
	}

	return nil
}

type connectionString struct {
	host     string
	port     int
	user     string
	password string
	database string
	sslMode  string
}

func newConnectionString(host string, user string, password string, database string, port int, sslMode string) connectionString {
	return connectionString{
		host:     host,
		user:     user,
		password: password,
		database: database,
		port:     port,
		sslMode:  sslMode,
	}
}

func (c connectionString) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		c.host, c.port, c.user, c.password, c.database, c.sslMode)

}

func getOrDefaultString(property string, defaultValue string) string {
	if viper.IsSet(property) {
		return viper.GetString(property)
	} else {
		return defaultValue
	}
}

func getOrDefaultInt(property string, defaultValue int) int {
	if viper.IsSet(property) {
		return viper.GetInt(property)
	} else {
		return defaultValue
	}
}
