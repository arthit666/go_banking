package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthit666/make_app/account"
	"github.com/arthit666/make_app/pocket"
	"github.com/arthit666/make_app/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// @title Banking API
// @version 1.0
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization

func main() {
	db, err := InitDb()
	if err != nil {
		panic("filed to connect to database")
	}

	db.AutoMigrate(
		&account.Account{},
		&account.AccountTransfer{},
		&pocket.Pocket{},
		&pocket.PocketTransfer{},
	)

	app := routes.RegRoute(db)

	go func() {
		if err := app.Listen(":8000"); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started on port 8000")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
	}
	log.Println("Server exited gracefully")
}
