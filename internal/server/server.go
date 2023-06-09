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
	port   string
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

	server := Server{server: serv, port: port}
	return &server, nil
}

func (serv *Server) Start() {
	log.Printf("Server running at port " + serv.port)
	log.Fatal(serv.server.ListenAndServe())
}
