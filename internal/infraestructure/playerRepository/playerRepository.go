package playerRepository

import (
	player "Lila-Back/pkg/Domain/player"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	GetPlayer(playerID int, txn *sql.Tx) (player.Player, int)
	CreatePlayer(player *player.Player, txn *sql.Tx) int
	UpdatePlayer(player player.Player, txn *sql.Tx) int
	DeletePlayer(playerID int, txn *sql.Tx) int
	DoMatchmaking(player player.Player, txn *sql.Tx) ([]player.Player, int)
}

type PlayerRepository struct{}

func (pr PlayerRepository) GetPlayer(playerID int, txn *sql.Tx) (player.Player, int) {
	var player player.Player
	stmt, err := txn.Prepare(`SELECT
								player_id,
								name,
								level,
								player_rank,
								winrate,
								active
							FROM 
								player
							WHERE
								player_id = ? AND
								active = True;`)
	defer stmt.Close()
	if err != nil {
		return player, http.StatusInternalServerError
	}

	err = stmt.QueryRow(playerID).Scan(
		&player.Id,
		&player.Name,
		&player.Level,
		&player.Rank,
		&player.Winrate,
		&player.Active)

	if err == sql.ErrNoRows {
		return player, http.StatusNotFound
	}
	if err != nil {
		return player, http.StatusInternalServerError
	}
	return player, http.StatusOK
}

func (pr PlayerRepository) CreatePlayer(player *player.Player, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`INSERT INTO player (name)
							VALUES (?);`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(player.Name)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "unique_name") {
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

func (pr PlayerRepository) UpdatePlayer(player player.Player, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`UPDATE 
								player 
							SET
								name = ?,
								level = ?,
								player_rank = ?,
								winrate = ?
							WHERE 
								player_id = ? AND
								active = True;`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(
		player.Name,
		player.Level,
		player.Rank,
		player.Winrate,
		player.Id,
	)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "unique_name") {
			return http.StatusBadRequest
		}
		if strings.Contains(str, "level_check") {
			return http.StatusBadRequest
		}
		if strings.Contains(str, "rank_check") {
			return http.StatusBadRequest
		}
		if strings.Contains(str, "winrate_check") {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return http.StatusNotModified
	}

	return http.StatusOK
}

func (pr PlayerRepository) DeletePlayer(playerID int, txn *sql.Tx) int {

	stmt, err := txn.Prepare(`UPDATE 
								player 
							SET
								active = False
							WHERE 
								player_id = ? AND
								active = true;`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(
		playerID,
	)
	if err != nil {
		return http.StatusInternalServerError
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return http.StatusNotFound
	}

	return http.StatusOK

}

func (pr PlayerRepository) DoMatchmaking(playerR player.Player, txn *sql.Tx) ([]player.Player, int) {
	stmt, err := txn.Prepare(`SELECT
								player_id,
								name,
								level,
								player_rank,
								winrate,
								active
							FROM 
								player
							WHERE 
								active = True
							ORDER BY
								ABS( ? - player_rank),
								ABS( ? - level ),
								ABS( ? - winrate )
							LIMIT 15;`)
	defer stmt.Close()
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	rows, err := stmt.Query(
		playerR.Rank,
		playerR.Level,
		playerR.Winrate,
	)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	var players []player.Player
	var noRows bool = true
	for rows.Next() {
		noRows = false
		var player player.Player
		err = rows.Scan(
			&player.Id,
			&player.Name,
			&player.Level,
			&player.Rank,
			&player.Winrate,
			&player.Active,
		)
		if err != nil {
			return nil, http.StatusInternalServerError
		}
		players = append(players, player)
	}

	if noRows {
		return nil, http.StatusNotFound
	}

	return players, http.StatusOK

}
