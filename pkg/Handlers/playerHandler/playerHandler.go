package playerHandler

import (
	"Lila-Back/internal/infraestructure/playerRepository"
	"Lila-Back/pkg/Domain/player"
	"database/sql"
	"net/http"
)

type PlayerHandler struct {
	Repository playerRepository.Repository
}

type Handler interface {
	GetPlayer(playerID int, txn *sql.Tx) (interface{}, int)
	CreatePlayer(player player.Player, txn *sql.Tx) int
	UpdatePlayer(player player.Player, txn *sql.Tx) int
	DeletePlayer(playerID int, txn *sql.Tx) int
	DoMatchmaking(player player.Player, txn *sql.Tx) (interface{}, int)
}

func (ph PlayerHandler) GetPlayer(playerID int, txn *sql.Tx) (interface{}, int) {
	player, status := ph.Repository.GetPlayer(playerID, txn)
	if status != http.StatusOK {
		return nil, status
	}
	return player, status
}

func (ph PlayerHandler) CreatePlayer(player player.Player, txn *sql.Tx) int {
	return ph.Repository.CreatePlayer(player, txn)
}

func (ph PlayerHandler) UpdatePlayer(player player.Player, txn *sql.Tx) int {
	return ph.Repository.UpdatePlayer(player, txn)
}

func (ph PlayerHandler) DeletePlayer(playerID int, txn *sql.Tx) int {
	return ph.Repository.DeletePlayer(playerID, txn)
}

func (ph PlayerHandler) DoMatchmaking(player player.Player, txn *sql.Tx) (interface{}, int) {
	players, status := ph.Repository.DoMatchmaking(player, txn)
	if status != http.StatusOK {
		return nil, status
	}
	// TODO retrun ch.R.DOMATCH
	return players, status
}
