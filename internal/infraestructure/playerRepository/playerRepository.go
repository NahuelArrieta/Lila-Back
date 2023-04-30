package playerRepository

import (
	player "Lila-Back/pkg/Domain/player"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	GetPlayer(playerID int, txn *sql.Tx) (player.Player, int)
	CreatePlayer(player player.Player, txn *sql.Tx) int
	// UpdatePlayer(playerID int, txn *sql.Tx) int
	// DeletePlayer(playerID int, txn *sql.Tx) int
}

type PlayerRepository struct{}

func (pr PlayerRepository) GetPlayer(playerID int, txn *sql.Tx) (player.Player, int) {
	var player player.Player
	stmt, err := txn.Prepare(`SELECT
								player_id,
								name,
								level,
								player_rank,
								winrate
							FROM player
							WHERE
								player_id = ?;`)
	if err != nil {
		return player, http.StatusInternalServerError
	}

	err = stmt.QueryRow(playerID).Scan(
		&player.Id,
		&player.Name,
		&player.Level,
		&player.Rank,
		&player.Winrate)

	if err == sql.ErrNoRows {
		return player, http.StatusNotFound
	}
	if err != nil {
		return player, http.StatusInternalServerError
	}
	return player, http.StatusOK
}

func (pr PlayerRepository) CreatePlayer(player player.Player, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`INSERT INTO player (
								name,
								level,
								player_rank,
								winrate)
							VALUES (?,?,?,?);`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(
		player.Name,
		player.Level,
		player.Rank,
		player.Winrate,
	)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "player.name") {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}

	id, err := res.LastInsertId()
	if err != nil {
		return http.StatusInternalServerError
	}
	player.Id = int(id)

	return http.StatusOK

}
