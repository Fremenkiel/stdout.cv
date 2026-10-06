package table

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type Handler struct {
	renderer	renderer
	service		*Service
}

func NewHandler(r renderer, s *Service) *Handler {
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

	viewData := h.getTableViewData(tables)

	if err := h.renderer.RenderPage(w, "index", viewData); err != nil {
		log.Printf("table: error thrown while rendering, %v", err)
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
		}
	}

	return viewData
}
