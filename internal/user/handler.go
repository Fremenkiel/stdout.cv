package user

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
	user, err := h.service.GetUser(r.Context())
	if err != nil {
		log.Printf("Error thrown: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		return
	}

	viewData := h.getUserViewData(user)

	if err := h.renderer.RenderPage(w, "user", viewData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) getUserViewData(user *User) UserViewData {
	return UserViewData{
		Name: user.Name,
	}
}
