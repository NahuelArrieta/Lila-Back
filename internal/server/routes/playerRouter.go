package routes

import (
	"Lila-Back/pkg/Domain/player"
	"Lila-Back/pkg/Domain/response"
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
		_, err = w.Write([]byte("Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	player, result := pr.Handler.GetPlayer(playerID, txn)
	w.WriteHeader(result.Status)
	resp, err := result.BuildResponse(player)
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
		result := response.Response{Message: "Bad Request"}
		resp, err := result.BuildResponse(nil)
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

	err = json.NewDecoder(r.Body).Decode(&player)
	if player.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		result := response.Response{Message: "Missing Parameters"}
		resp, err := result.BuildResponse(nil)
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

	player.SetDefaultValues()

	result := pr.Handler.CreatePlayer(&player, txn)

	w.WriteHeader(result.Status)

	resp, err := result.BuildResponse(player)
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

	defer txn.Rollback()
	if result.Status == http.StatusOK {
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

	player := player.Player{Id: -1}

	err = json.NewDecoder(r.Body).Decode(&player)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := response.Response{Message: "Bad Request"}
		resp, err := result.BuildResponse(nil)
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

	if player.Id == -1 || player.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		result := response.Response{Message: "Missing parameters"}
		resp, err := result.BuildResponse(nil)
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

	// Not allowed to change the value active
	player.Active = true

	result := pr.Handler.UpdatePlayer(player, txn)
	w.WriteHeader(result.Status)

	resp, err := result.BuildResponse(player)
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

	defer txn.Rollback()
	if result.Status == http.StatusOK {
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
		result := response.Response{Message: "Bad Request"}
		resp, err := result.BuildResponse(nil)
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

	result := pr.Handler.DeletePlayer(playerID, txn)
	w.WriteHeader(result.Status)

	defer txn.Rollback()
	if result.Status == http.StatusOK {
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
	playerStats := player.PlayerStats{Level: -1, Rank: -1, Winrate: -1}

	err = json.NewDecoder(r.Body).Decode(&playerStats)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := response.Response{Message: "Bad Request"}
		resp, err := result.BuildResponse(nil)
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

	if playerStats.Level == -1 || playerStats.Rank == -1 || playerStats.Winrate == -1 {
		w.WriteHeader(http.StatusBadRequest)
		result := response.Response{Message: "Missing Parameters"}
		resp, err := result.BuildResponse(nil)
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

	players, result := pr.Handler.DoMatchmaking(playerStats, txn)
	w.WriteHeader(result.Status)
	resp, err := result.BuildResponse(players)
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
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	r.Get("/{playerId}", pr.GetPlayer)
	r.Post("/", pr.CreatePlayer)
	r.Put("/", pr.UpdatePlayer)
	r.Delete("/{playerId}", pr.DeletePlayer)

	r.Put("/matchmaking", pr.DoMatchmaking)

	return r
}
