package playerRepository

import (
	player "Lila-Back/pkg/Domain/player"
	"database/sql"
)

type Repository interface {
	GetPlayer(playerID int, txn *sql.Tx) (player.Player, int)
	CreatePlayer(player player.Player, txn *sql.Tx) int
	UpdatePlayer(playerID int, txn *sql.Tx) int
	DeletePlayer(playerID int, txn *sql.Tx) int
}

type playerRepository struct{}

func GetPlayer(playerID int, txn *sql.Tx) (player.Player, int) {
	var player player.Player
	stmt, err := txn.Prepare(`SELECT 
								player_id,
								name,
								level,
								player_rank,
								winrate
							FROM player
							WHERE
								player_id = ?`)
	if err != nil {
		return player, 0 //TODO arreglar
	}

	err = stmt.QueryRow(playerID).Scan(
		&player.Id,
		&player.Name,
		&player.Level,
		&player.Rank,
		&player.Winrate)

	// TODO Reconocer error
	if err != nil {
		return player, -1
	}
	return player, 300 //TODO status
}
