package row

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/platform/lib"
	"github.com/fremenkiel/stdout.cv/internal/query"
	"github.com/fremenkiel/stdout.cv/internal/ui/viewmodels"
)

type Renderer interface {
	RenderPage(w http.ResponseWriter, name string, data any) error
	RenderFragment(w http.ResponseWriter, name string, data any) error
}

type Handler struct {
	renderer	Renderer
	service		*Service
}

func NewHandler(r Renderer, s *Service) *Handler {
	return &Handler{renderer: r, service: s}
}

func (h *Handler) GetRows(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryString := r.URL.Query().Get("query")
	if len(queryString) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := h.service.GetRows(ctx, queryString)
	if err != nil {
		if errors.Is(err, query.ErrQueryNotAllowed) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			libErr := &lib.Error{
				Title: "Query not allowed",
				Message: "The entered query is not allowed at this point in time",
			}
			if err = h.renderer.RenderFragment(w, "error-container-swap", libErr); err != nil {
				log.Printf("row: error thrown while rendering error container, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if libErr, ok := lib.ParseSqliteError(err); ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err = h.renderer.RenderFragment(w, "error-container-swap", libErr); err != nil {
				log.Printf("row: error thrown while rendering error container, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		log.Printf("Error thrown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	if ctx.Err() != nil {
		log.Printf("row: context error, %v", err)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	viewData := h.getResultViewData(result)

	if err := h.renderer.RenderFragment(w, "result", viewData); err != nil {
		log.Printf("row: error thrown while rendering, %v", err)
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
