package server

import (
	"Lila-Back/internal/server/routes"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	server *http.Server
}

func New(port string) (*Server, error) {
	r := chi.NewRouter()

	r.Mount("/api/videogame", routes.New())

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
	}

	server := Server{server: serv}
	return &server, nil
}

func (serv *Server) Start() {
	log.Printf("Servidor corriendo")
	log.Fatal(serv.server.ListenAndServe())
}
