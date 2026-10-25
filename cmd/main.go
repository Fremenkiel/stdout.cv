package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/platform/middleware"
	"github.com/fremenkiel/stdout.cv/internal/query"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/table"
	"github.com/fremenkiel/stdout.cv/internal/ui"
	"github.com/fremenkiel/stdout.cv/pkg/dotenv"
)

func main() {
	logger.NewLogger()

	if err := dotenv.Load(); err != nil {
		logger.LogFatal(
			"main",
			"unable to boot",
			slog.Any("error", err),
			)
	}

	_, err := os.Stat("/tmp/stdout_cv_sessions")
	if err != nil {
		if err := os.MkdirAll("/tmp/stdout_cv_sessions", os.ModePerm); err != nil {
			logger.LogFatal(
				"main",
				"unable to create tmp folder",
				slog.Any("error", err),
				)
		}
	}

	sessionCache := session.NewCache(
		os.Getenv("DATABASE_PATH"),
		os.Getenv("DATABASE_FILE_TEMPLATE"),
		)

	queryRepo := query.NewRepository(sessionCache)
	schemaRepo := schema.NewRepository(sessionCache)
	tableRepo := table.NewRepository(sessionCache)

	schemaService := schema.NewService(schemaRepo)
	rowService := query.NewService(queryRepo, schemaService)
	sessionService := session.NewService(sessionCache)
	tableService := table.NewService(tableRepo, schemaService)

	renderer := ui.NewRenderer()

	rowHandler := query.NewHandler(renderer, rowService)
	tableHandler := table.NewHandler(renderer, tableService)

	mux := http.NewServeMux()

	sessionMiddleware := session.Middleware(sessionService)

	mux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./ui/static/scripts"))))
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./ui/static/styles"))))

	requireHtml := middleware.EnforceAccept("text/html")

	table.NewRouter(
		mux,
		tableHandler,
		logger.LogRequestHandler,
		requireHtml,
		sessionMiddleware,
		)
	query.NewRouter(
		mux,
		rowHandler,
		logger.LogRequestHandler,
		requireHtml,
		sessionMiddleware,
		)

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)

	logger.LogInfo(
		"main",
		"server started",
		slog.String("address", "127.0.0.1"),
		slog.String("port", port),
		)

	server := http.Server{
		Addr: address,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sessionService.StartCleanupWorker(ctx, time.Duration(5 * time.Minute))

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.LogInfo(
				"main",
				"http listen and serve error",
				slog.Any("error", err),
				)
		}
	}()

	<-ctx.Done()
	startShutdown := time.Now()
	logger.LogInfo(
		"main",
		"shutdown signel recived",
		slog.Int64("recived_at", startShutdown.Unix()),
		)

	shotdownCtx, shotdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shotdownRelease()

	if err := server.Shutdown(shotdownCtx); err != nil {
		logger.LogInfo(
			"main",
			"unable to gracefully shutdown http server",
			slog.Any("error", err),
			)
	}

	logger.LogInfo(
		"main",
		"shutdown complete",
		slog.Int64("duration_ms", time.Now().Sub(startShutdown).Milliseconds()),
		)
}
