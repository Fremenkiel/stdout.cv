package row

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/fremenkiel/stdout.cv/internal/interpreter"
	"github.com/fremenkiel/stdout.cv/internal/platform/lib"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type RowService interface {
	GetRows(ctx context.Context, query string) (*Result, error)
}

type Handler struct {
	renderer	Renderer
	service		RowService
}

func NewHandler(r Renderer, s RowService) *Handler {
	return &Handler{renderer: r, service: s}
}

func (h *Handler) GetRows(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query().Get("query")
	if len(query) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := h.service.GetRows(ctx, query)
	if err != nil {
		if typedErr, ok := errors.AsType[*interpreter.Error](err); ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			libErr := &lib.Error{
				Title: "Error",
				Message: typedErr.Error(),
				Code: int(typedErr.Code()),
			}
			if err = h.renderer.RenderFragment(w, "error-container-swap", libErr); err != nil {
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
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err = h.renderer.RenderFragment(w, "error-container-swap", libErr); err != nil {
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
		w.WriteHeader(http.StatusInternalServerError)
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

	viewData := h.getResultViewData(result)

	if err := h.renderer.RenderFragment(w, "result", viewData); err != nil {
		logger.LogError(
			"row/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getResultViewData(result *Result) *viewmodels.ResultViewData {
	viewData := &viewmodels.ResultViewData{
		Stats: &viewmodels.Stats{
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
