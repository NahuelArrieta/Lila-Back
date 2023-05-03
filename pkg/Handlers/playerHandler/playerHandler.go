package playerHandler

import (
	"Lila-Back/internal/infraestructure/playerRepository"
	"Lila-Back/pkg/Domain/player"
	"Lila-Back/pkg/Domain/response"
	"database/sql"
	"net/http"
)

type PlayerHandler struct {
	Repository playerRepository.Repository
}

type Handler interface {
	GetPlayer(playerID int, txn *sql.Tx) (interface{}, response.Response)
	CreatePlayer(player *player.Player, txn *sql.Tx) response.Response
	UpdatePlayer(player player.Player, txn *sql.Tx) response.Response
	DeletePlayer(playerID int, txn *sql.Tx) response.Response
	DoMatchmaking(player player.PlayerStats, txn *sql.Tx) (interface{}, response.Response)
}

func (ph PlayerHandler) GetPlayer(playerID int, txn *sql.Tx) (interface{}, response.Response) {
	player, resp := ph.Repository.GetPlayer(playerID, txn)
	if resp.Status != http.StatusOK {
		return nil, resp
	}
	return player, resp
}

func (ph PlayerHandler) CreatePlayer(player *player.Player, txn *sql.Tx) response.Response {
	return ph.Repository.CreatePlayer(player, txn)
}

func (ph PlayerHandler) UpdatePlayer(player player.Player, txn *sql.Tx) response.Response {
	_, resp := ph.Repository.GetPlayer(player.Id, txn)
	// Verify wether the player exists
	if resp.Status == http.StatusNotFound {
		return resp
	}
	return ph.Repository.UpdatePlayer(player, txn)
}

func (ph PlayerHandler) DeletePlayer(playerID int, txn *sql.Tx) response.Response {
	return ph.Repository.DeletePlayer(playerID, txn)
}

func (ph PlayerHandler) DoMatchmaking(playerStats player.PlayerStats, txn *sql.Tx) (interface{}, response.Response) {
	return ph.Repository.DoMatchmaking(playerStats, txn)
}
