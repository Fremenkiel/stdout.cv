package row

import (
	"encoding/json"
	"errors"
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

func (h *Handler) GetRows(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query().Get("query")

	rows, err := h.service.GetRows(ctx, query)
	if err != nil {
		if errors.Is(err, ErrQueryNotAllowed) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		log.Printf("Error thrown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
		return
	}

	if ctx.Err() != nil {
		log.Printf("row: context error, %v", err)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	viewData := &RowListViewData{}

	for i := range rows {
		viewData.Rows = append(viewData.Rows, h.getRowViewData(rows[i]))
	}

	if err := h.renderer.RenderFragment(w, "row-container", viewData); err != nil {
		log.Printf("row: error thrown while rendering, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getRowViewData(row *Row) RowViewData {
	var viewData RowViewData

	for _, str := range *row {
		viewData = append(viewData, *str)
	}

	return viewData
}
