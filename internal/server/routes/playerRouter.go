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
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	player, status := pr.Handler.GetPlayer(playerID, txn) //TODO recuperar status
	w.WriteHeader(status)
	resp, _ := json.Marshal(player)
	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}
}

func (pr playerRouter) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	//TODO
	txn, _ := connection.Connect()
	player := &player.Player{
		Name:    "Player0",
		Level:   "0",
		Rank:    "0",
		Winrate: "0",
	}

	status := pr.Handler.CreatePlayer(*player, txn)
	w.WriteHeader(status)

}

func (pr *playerRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/{playerId}", pr.GetPlayer)
	r.Post("/", pr.CreatePlayer)

	return r
}
