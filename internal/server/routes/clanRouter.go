package routes

import (
	"Lila-Back/pkg/Domain/clan"
	"Lila-Back/pkg/Handlers/clanHandler"
	connection "Lila-Back/pkg/Helpers/Connection"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type clanRouter struct {
	Handler clanHandler.Handler
}

func (cr clanRouter) CreateClan(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	clan := clan.Clan{LeaderId: -1}
	err = json.NewDecoder(r.Body).Decode(&clan)
	if err != nil || clan.Name == "" || clan.HashedPassword == "" || clan.LeaderId == -1 {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	status := cr.Handler.CreateClan(&clan, txn)

	resp, err := json.Marshal(clan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	w.WriteHeader(status)
	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}
}

func (cr clanRouter) JoinClan(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	jr := clan.JoinRequest{Clan_Id: -1, Player_Id: -1}
	err = json.NewDecoder(r.Body).Decode(&jr)
	if err != nil || jr.Clan_Id == -1 || jr.Player_Id == -1 {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	status := cr.Handler.JoinClan(&jr, txn)

	resp, err := json.Marshal(jr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	w.WriteHeader(status)
	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}
}

func (cr clanRouter) PutColeader(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	coleaderReq := clan.ColeaderRequest{ClanId: -1, LeaderId: -1, ColeaderId: -1}
	err = json.NewDecoder(r.Body).Decode(&coleaderReq)
	if err != nil || coleaderReq.ClanId == -1 || coleaderReq.LeaderId == -1 || coleaderReq.ColeaderId == -1 {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	status := cr.Handler.PutColeader(coleaderReq, txn)

	resp, err := json.Marshal(coleaderReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	w.WriteHeader(status)
	_, err = w.Write([]byte(resp))
	if err != nil {
		fmt.Println("Internal Fatal Error")
	}

	defer txn.Rollback()
	if status == http.StatusOK {
		txn.Commit()
	}
}

func (cr clanRouter) GetPlayers(w http.ResponseWriter, r *http.Request) {
	txn, err := connection.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("500: Internal Server Error"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	clanID, err := strconv.Atoi(chi.URLParam(r, "clanId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// _, err = w.Write([]byte("400 Bad Request"))
		if err != nil {
			fmt.Println("Internal Fatal Error")
		}
		return
	}

	players, status := cr.Handler.GetPlayers(clanID, txn)
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

func (cr clanRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"}, //--
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	r.Post("/", cr.CreateClan)
	r.Put("/join", cr.JoinClan)
	r.Put("/co-leader", cr.PutColeader)
	r.Get("/players/{clanId}", cr.GetPlayers)

	return r

}
