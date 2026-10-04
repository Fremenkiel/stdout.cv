package table

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/platform/render"
)

type Handler struct {
	renderer	render.Renderer
	service		*Service
}

func NewHandler(r render.Renderer, s *Service) *Handler {
	return &Handler{renderer: r, service: s}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tables, err := h.service.GetTables(ctx)
	if err != nil {
		log.Printf("Error thrown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tables)
		return
	}

	if ctx.Err() != nil {
		log.Printf("table: context error, %v", ctx.Err())
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	viewData := &TableListViewData{
		Tables: make([]TableViewData, len(tables)),
		Error: nil,
		Rows: nil,
	}

	for i := range tables {
		viewData.Tables[i] = h.getTableViewData(tables[i])
	}

	if err := h.renderer.RenderPage(w, "index", viewData); err != nil {
		log.Printf("table: error thrown while rendering, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getTableViewData(table *Table) TableViewData {
	return TableViewData{
		Name: table.Name,
	}
}
