package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/platform/render"
	"github.com/fremenkiel/stdout.cv/internal/user"
	"github.com/fremenkiel/stdout.cv/pkg/dotenv"
	_ "modernc.org/sqlite"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env files: %v", err)
	}

	dbURI := os.Getenv("DB_URI")

	db, err := sql.Open("sqlite", dbURI)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v", err)
	}

	if err := database.InitializeSchema(db); err != nil {
		log.Fatal("Unable to apply schema: %v", err)
	}

	mux := &http.ServeMux{}

	mux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./ui/static/scripts"))))

	renderer := render.NewTemplateRenderer()

	userRepo := user.NewRepository(db)

	userService := user.NewService(userRepo)

	userHandler := user.NewHandler(renderer, userService)

	user.NewRouter(mux, userHandler)

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	log.Printf("Server listing on 127.0.0.1%s", address)

	server := http.Server{
		Addr: address,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shotdownCtx, shotdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shotdownRelease()

	if err := server.Shutdown(shotdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Print("Shutdown complete")
}
