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
	// UpdatePlayer(playerID int, txn *sql.Tx) int
	// DeletePlayer(playerID int, txn *sql.Tx) int
}

func (ph PlayerHandler) GetPlayer(playerID int, txn *sql.Tx) (interface{}, int) {
	player, status := ph.Repository.GetPlayer(playerID, txn)
	if status != http.StatusOK {
		return nil, status
	}
	return player, status
}

func (ph PlayerHandler) CreatePlayer(player player.Player, txn *sql.Tx) int {
	// TODO
	return ph.Repository.CreatePlayer(player, txn)
}
