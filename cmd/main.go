package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/middleware"
	"github.com/fremenkiel/stdout.cv/internal/platform/render"
	"github.com/fremenkiel/stdout.cv/internal/row"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/table"
	"github.com/fremenkiel/stdout.cv/pkg/dotenv"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Fatalf("Unable to load .env files: %v", err)
	}

	_, err := os.Stat("/tmp/stdout_cv_sessions")
	if err != nil {
		if err := os.MkdirAll("/tmp/stdout_cv_sessions", os.ModePerm); err != nil {
			log.Fatalf("Unable to create tmp folder: %v", err)
		}
	}

	sessionCache := session.NewCache()

	rowRepo := row.NewRepository(sessionCache)
	schemaRepo := schema.NewRepository(sessionCache)
	tableRepo := table.NewRepository(sessionCache)

	schemaService := schema.NewService(schemaRepo)
	rowService := row.NewService(rowRepo, schemaService)
	sessionService := session.NewService(sessionCache)
	tableService := table.NewService(tableRepo)

	renderer := render.NewTemplateRenderer()

	rowHandler := row.NewHandler(renderer, rowService)
	tableHandler := table.NewHandler(renderer, tableService)

	mux := new(middleware.MiddlewareMux)

	sessionMiddleware := session.NewMiddleware(sessionService)
	mux.AppendMiddleware(sessionMiddleware.Handle)

	mux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./ui/static/scripts"))))
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./ui/static/styles"))))

	table.NewRouter(mux, tableHandler)
	row.NewRouter(mux, rowHandler)

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
