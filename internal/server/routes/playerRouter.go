package routes

import (
	"Lila-Back/pkg/Domain/player"
	"Lila-Back/pkg/Handlers/playerHandler"
	connection "Lila-Back/pkg/Helpers/Connection"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type playerRouter struct {
	Handler playerHandler.Handler
}

// var qlCon *sql.DB  --- TODO

func (pr playerRouter) GetPlayer(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	player, status := pr.Handler.GetPlayer(playerID, txn) //TODO recuperar status
	w.WriteHeader(status)
	resp, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}
}

func (pr playerRouter) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	var player player.Player

	err = json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	player.Level = 0
	player.Rank = 0
	player.Winrate = 0
	player.Active = true

	status := pr.Handler.CreatePlayer(player, txn)
	w.WriteHeader(status)

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}

}

func (pr playerRouter) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	var player player.Player

	err = json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}
	player.Id = playerID

	status := pr.Handler.UpdatePlayer(player, txn) // TODO cambia
	w.WriteHeader(status)

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}

}

func (pr playerRouter) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	playerID, err := strconv.Atoi(chi.URLParam(r, "playerId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	status := pr.Handler.DeletePlayer(playerID, txn) // TODO cambia
	w.WriteHeader(status)

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}
}

func (pr playerRouter) DoMatchmaking(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}
	var player player.Player

	err = json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	players, status := pr.Handler.DoMatchmaking(player, txn)
	w.WriteHeader(status)
	resp, err := json.Marshal(players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}

}

func (pr *playerRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, //--
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	r.Get("/{playerId}", pr.GetPlayer)
	r.Post("/", pr.CreatePlayer)
	r.Put("/{playerId}", pr.UpdatePlayer)
	r.Delete("/{playerId}", pr.DeletePlayer)

	r.Put("/matchmaking", pr.DoMatchmaking)

	return r
}
