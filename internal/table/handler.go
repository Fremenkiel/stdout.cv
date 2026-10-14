package table

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type TableService interface {
	GetTables(ctx context.Context) ([]*Table, error)
}

type Handler struct {
	renderer	Renderer
	service		TableService
}

func NewHandler(r Renderer, s TableService) *Handler {
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

	if err := ctx.Err(); err != nil {
		logger.LogError(
			"table/handler",
			"context error",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	viewData := h.getTableViewData(tables)

	if err := h.renderer.RenderPage(w, "index", viewData); err != nil {
		logger.LogError(
			"table/handler",
			"error returned while rendering",
			slog.Any("error", err),
			)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getTableViewData(tables []*Table) *viewmodels.TableViewData {
	viewData := &viewmodels.TableViewData{
		Tables: make([]viewmodels.Table, len(tables)),
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
