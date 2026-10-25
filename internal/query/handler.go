package query

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/interpreter"
	"github.com/fremenkiel/stdout.cv/internal/platform/lib"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type service interface {
	Run(ctx context.Context, sql string) (*Query, error)
}

type Handler struct {
	renderer	Renderer
	service		service
}

func NewHandler(r Renderer, s service) *Handler {
	return &Handler{renderer: r, service: s}
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
			libErr.Duration =time.Now().UnixMicro() - start.UnixMicro()
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

	result.Duration = time.Now().UnixMicro() - start.UnixMicro()
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

func (h *Handler) getResultViewData(result *Query) *viewmodels.ResultViewData {
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
