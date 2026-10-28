package query

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/database"
	"github.com/fremenkiel/stdout.cv/internal/interpreter"
	"github.com/fremenkiel/stdout.cv/internal/platform/lib"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type service interface {
	Run(ctx context.Context, sql string) (*Query, error)
	GetSchema(ctx context.Context) (schema.Schema, error)
	GetStats(ctx context.Context) (*database.Stats, error)
}

type Handler struct {
	renderer	renderer
	service		service
}

func NewHandler(r renderer, s service) *Handler {
	return &Handler{renderer: r, service: s}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	schema, err := h.service.GetSchema(ctx)
	if err != nil {
		logger.LogError(
			"table/handler",
			"error returned while getting tables",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats, err := h.service.GetStats(ctx)
	if err != nil {
		logger.LogError(
			"table/handler",
			"error returned while getting stats",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ctx.Err(); err != nil {
		logger.LogError(
			"table/handler",
			"context error",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	viewData := h.getViewData(&Query{
		Schema: schema,
		Stats: stats,
		Result: nil,
	}, nil)

	if err := h.renderer.RenderPage(w, "index", viewData); err != nil {
		logger.LogError(
			"table/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sql := r.URL.Query().Get("sql")
	if len(sql) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	start := time.Now()

	result, err := h.service.Run(ctx, sql)
	if err != nil {
		if typedErr, ok := errors.AsType[*interpreter.Error](err); ok {
			ln, csn, cen := typedErr.Pos()
			libErr := &lib.Error{
				Message: typedErr.Error(),
				Code: interpreter.ErrorCodeString[typedErr.Code()],
				Duration: time.Now().UnixMicro() - start.UnixMicro(),
				LineNumber: ln,
				CharStartNumber: csn,
				CharEndNumber: cen,
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err = h.renderer.RenderFragment(w, "error-result-swap", libErr); err != nil {
				logger.LogError(
					"row/handler",
					"error returned while rendering error container",
					slog.Any("error", err),
					)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if libErr, ok := lib.ParseSqliteError(err); ok {
			libErr.Duration = time.Now().UnixMicro() - start.UnixMicro()
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err = h.renderer.RenderFragment(w, "error-result-swap", libErr); err != nil {
				logger.LogError(
					"row/handler",
					"error returned while rendering error container",
					slog.Any("error", err),
					)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		logger.LogError(
			"row/handler",
			"error returned while getting rows",
			slog.Any("error", err),
			)
		libErr := &lib.Error{
			Message: err.Error(),
			Code: "INTERNAL_SERVER_ERROR",
			Duration: time.Now().UnixMicro() - start.UnixMicro(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err = h.renderer.RenderFragment(w, "error-result-swap", libErr); err != nil {
			logger.LogError(
				"row/handler",
				"error returned while rendering error container",
				slog.Any("error", err),
				)
		}
		return
	}

	if err := ctx.Err(); err != nil {
		logger.LogError(
			"row/handler",
			"context error",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	if result.Result != nil {
		result.Result.Duration = time.Now().UnixMicro() - start.UnixMicro()
	}

	viewData := h.getViewData(result, nil)

	if err := h.renderer.RenderFragment(w, "result-swap", viewData); err != nil {
		logger.LogError(
			"row/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getViewData(result *Query, err *lib.Error) *viewmodels.ViewData {
	viewData := &viewmodels.ViewData{
		Tables: h.getTableViewData(result.Schema),
		Stats: h.getStatsViewData(result.Stats),
		Error: err,
		Result: h.getResultViewData(result.Result),
	}

	return viewData
}

func (h *Handler) getTableViewData(tables schema.Schema) []*viewmodels.Table {
	if tables == nil {
		return nil
	}

	var viewData = make([]*viewmodels.Table, len(tables))

	for i := range tables {
		viewData[i] = &viewmodels.Table{
			Name: tables[i].Name,
			RowCount: tables[i].RowCount,
			Columns: h.getColumnViewData(tables[i].Columns),
		}
	}

	return viewData
}

func (h *Handler) getColumnViewData(columns []*schema.Column) *viewmodels.ColumnList {
	viewData := make(viewmodels.ColumnList, len(columns))

	for i := range columns {
		viewData[i] = viewmodels.Column{
			Name: columns[i].Name,
			Type: columns[i].Type,
			IsPrimaryKey: columns[i].IsPrimaryKey,
			IsForeignKey: columns[i].IsForeignKey,
			IsNullable: columns[i].IsNullable,
		}
	}

	return &viewData
}

func (h *Handler) getStatsViewData(stats *database.Stats) *viewmodels.Stats {
	if stats == nil {
		return nil
	}

	return &viewmodels.Stats{
		Version: stats.Version,
		Size: stats.Size,
	}
}

func (h *Handler) getResultViewData(result *Result) *viewmodels.Result {
	if result == nil {
		return nil
	}

	viewData := &viewmodels.Result{
		Stats: &viewmodels.ResultStats{
			RowCount: uint16(len(result.Rows)),
			ColumnCount: uint16(len(result.Columns)),
			Duration: result.Duration,
		},
		Rows: make([]viewmodels.Row, len(result.Rows)),
		Columns: make(viewmodels.ColumnList, len(result.Columns)),
	}

	for i := range result.Rows {
		viewData.Rows[i] = h.getRowViewModel(result.Rows[i])
	}

	for i, column := range result.Columns {
		viewData.Columns[i].Name = column
	}

	return viewData
}

func (h *Handler) getRowViewModel(row *Row) viewmodels.Row {
	viewModel := make(viewmodels.Row, len(*row))

	for i, val := range *row {
		viewModel[i] = *val
	}

	return viewModel
}
