package playerRepository

import (
	player "Lila-Back/pkg/Domain/player"
	"Lila-Back/pkg/Domain/response"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	GetPlayer(playerID int, txn *sql.Tx) (player.Player, response.Response)
	CreatePlayer(player *player.Player, txn *sql.Tx) response.Response
	UpdatePlayer(player player.Player, txn *sql.Tx) response.Response
	DeletePlayer(playerID int, txn *sql.Tx) response.Response
	DoMatchmaking(player player.PlayerStats, txn *sql.Tx) ([]player.Player, response.Response)
}

type PlayerRepository struct{}

func (pr PlayerRepository) GetPlayer(playerID int, txn *sql.Tx) (player.Player, response.Response) {
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
		return player, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	err = stmt.QueryRow(playerID).Scan(
		&player.Id,
		&player.Name,
		&player.Level,
		&player.Rank,
		&player.Winrate,
		&player.Active)

	if err == sql.ErrNoRows {
		return player, response.Response{Status: http.StatusNotFound, Message: "Player Not Found"}
	}
	if err != nil {
		return player, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	return player, response.Response{Status: http.StatusOK, Message: "Ok"}
}

func (pr PlayerRepository) CreatePlayer(player *player.Player, txn *sql.Tx) response.Response {
	stmt, err := txn.Prepare(`INSERT INTO player (name)
							VALUES (?);`)
	defer stmt.Close()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	res, err := stmt.Exec(player.Name)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "unique_name") {
			return response.Response{Status: http.StatusBadRequest, Message: "Name Already Used"}
		}
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	player.Id = int(id)

	return response.Response{Status: http.StatusOK, Message: "Player Created"}

}

func (pr PlayerRepository) UpdatePlayer(player player.Player, txn *sql.Tx) response.Response {
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
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	_, err = stmt.Exec(
		player.Name,
		player.Level,
		player.Rank,
		player.Winrate,
		player.Id,
	)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "unique_name") {
			return response.Response{Status: http.StatusBadRequest, Message: "Name Already Used"}
		}
		if strings.Contains(str, "level_check") {
			return response.Response{Status: http.StatusBadRequest, Message: "Level Can't be less than 0"}
		}
		if strings.Contains(str, "rank_check") {
			return response.Response{Status: http.StatusBadRequest, Message: "Rank Can't be less than 0"}
		}
		if strings.Contains(str, "winrate_check") {
			return response.Response{Status: http.StatusBadRequest, Message: "Winrate Must be between 0 and 100"}
		}
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	return response.Response{Status: http.StatusOK, Message: "Player Updated"}
}

func (pr PlayerRepository) DeletePlayer(playerID int, txn *sql.Tx) response.Response {

	stmt, err := txn.Prepare(`UPDATE 
								player 
							SET
								active = False
							WHERE 
								player_id = ? AND
								active = true;`)
	defer stmt.Close()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	res, err := stmt.Exec(
		playerID,
	)
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return response.Response{Status: http.StatusNotFound, Message: "Player Not Found"}
	}

	return response.Response{Status: http.StatusOK, Message: "Player Deleted"}

}

func (pr PlayerRepository) DoMatchmaking(playerStats player.PlayerStats, txn *sql.Tx) ([]player.Player, response.Response) {
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
		return nil, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	rows, err := stmt.Query(
		playerStats.Rank,
		playerStats.Level,
		playerStats.Winrate,
	)
	if err != nil {
		return nil, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var players []player.Player
	for rows.Next() {
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
			return nil, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
		}
		players = append(players, player)
	}

	return players, response.Response{Status: http.StatusOK, Message: "OK"}

}
