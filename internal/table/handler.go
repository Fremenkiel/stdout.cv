package table

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/database"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type service interface {
	GetTables(ctx context.Context) ([]*Table, error)
	GetStats(ctx context.Context) (*database.Stats, error)
}

type Handler struct {
	renderer	Renderer
	service		service
}

func NewHandler(r Renderer, s service) *Handler {
	return &Handler{renderer: r, service: s}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tables, err := h.service.GetTables(ctx)
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

	viewData := h.getTableViewData(tables, stats)

	if err := h.renderer.RenderPage(w, "index", viewData); err != nil {
		logger.LogError(
			"table/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getTableViewData(tables []*Table, stats *database.Stats) *viewmodels.TableViewData {
	viewData := &viewmodels.TableViewData{
		Tables: make([]viewmodels.Table, len(tables)),
		Stats: viewmodels.DatabaseStats{
			Version: stats.Version,
			Size: stats.Size,
		},
		Error: nil,
		Result: &viewmodels.Result{},
	}

	for i := range tables {
		viewData.Tables[i] = viewmodels.Table{
			Name: tables[i].Name,
			RowCount: tables[i].RowCount,
			Columns: h.getColumnViewData(tables[i].Columns),
		}
	}

	return viewData
}

func (h *Handler) getColumnViewData(columns []*column.Column) *viewmodels.ColumnList {
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
