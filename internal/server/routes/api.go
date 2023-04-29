package routes

import (
	"Lila-Back/internal/infraestructure/playerRepository"
	"Lila-Back/pkg/Handlers/playerHandler"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
)

var (
	sqlCon *sql.DB
)

func New() http.Handler {
	r := chi.NewRouter()

	pr := &playerRouter{
		Handler: &playerHandler.PlayerHandler{
			Repository: &playerRepository.PlayerRepository{},
		},
	}

	r.Mount("/player", pr.Routes()) //

	return r
}
