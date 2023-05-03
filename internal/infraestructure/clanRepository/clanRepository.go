package clanRepository

import (
	"Lila-Back/pkg/Domain/clan"
	"Lila-Back/pkg/Domain/player"
	"Lila-Back/pkg/Domain/response"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) response.Response
	JoinClan(playerID, clanID int, txn *sql.Tx) response.Response
	GetClanPassword(clanID int, txn *sql.Tx) (string, response.Response)
	PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) response.Response
	GetClanLeader(clanID int, txn *sql.Tx) (int, response.Response)
	GetPlayers(clanID int, txn *sql.Tx) ([]player.Player, response.Response)
}

type ClanRepository struct{}

func (cr ClanRepository) CreateClan(clan *clan.Clan, txn *sql.Tx) response.Response {
	stmt, err := txn.Prepare(`INSERT INTO clan (
								name,
								hashed_password,
								leader_id)
							VALUES (?,?,?);`)
	defer stmt.Close()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	res, err := stmt.Exec(
		clan.Name,
		clan.HashedPassword,
		clan.LeaderId,
	)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "unique_name") {
			return response.Response{Status: http.StatusBadRequest, Message: "Name Already Used"}
		}
		if strings.Contains(str, "leader_fk") {
			return response.Response{Status: http.StatusNotFound, Message: "Leader Not Found"}
		}
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	clan.Id = int(id)

	return response.Response{Status: http.StatusOK, Message: "Clan Created"}
}

func (cr ClanRepository) JoinClan(playerID int, clanID int, txn *sql.Tx) response.Response {
	stmt, err := txn.Prepare(`INSERT INTO clan_player(
								player_id,
								clan_id)
							VALUES (?,?);`)
	defer stmt.Close()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	_, err = stmt.Exec(playerID, clanID)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "Duplicate") {
			return response.Response{Status: http.StatusOK, Message: "Player is Already in the Clan"}
		}
		if strings.Contains(str, "player_fk") {
			return response.Response{Status: http.StatusNotFound, Message: "Player Not Found"}
		}
		if strings.Contains(str, "clan_fk") {
			return response.Response{Status: http.StatusNotFound, Message: "Clan Not Found"}
		}
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	return response.Response{Status: http.StatusOK, Message: "Player Joined to the Clan"}

}

func (cr ClanRepository) GetClanPassword(clanID int, txn *sql.Tx) (string, response.Response) {
	stmt, err := txn.Prepare(`SELECT 
								hashed_password
							FROM 
								clan
							WHERE 
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return "", response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var password string
	err = stmt.QueryRow(clanID).Scan(&password)

	if err == sql.ErrNoRows {
		return "", response.Response{Status: http.StatusNotFound, Message: "Clan Not Found"}
	}
	if err != nil {
		return "", response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	return password, response.Response{Status: http.StatusOK, Message: "OK"}

}

func (cr ClanRepository) PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) response.Response {
	stmt, err := txn.Prepare(`UPDATE 
								clan
							SET 
								coleader_id = ?
							WHERE
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	_, err = stmt.Exec(coleaderReq.ColeaderId, coleaderReq.ClanId)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "coleader_fk") {
			return response.Response{Status: http.StatusNotFound, Message: "Coleader Not Found"}
		}
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	return response.Response{Status: http.StatusOK, Message: "Coleader Asigned"}
}

func (cr ClanRepository) GetClanLeader(clanID int, txn *sql.Tx) (int, response.Response) {
	stmt, err := txn.Prepare(`SELECT 
								leader_id
							FROM 
								clan
							WHERE 
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return -1, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var LeaderId int
	err = stmt.QueryRow(clanID).Scan(&LeaderId)

	if err == sql.ErrNoRows {
		return -1, response.Response{Status: http.StatusNotFound, Message: "Clan Not Found"}
	}
	if err != nil {
		return -1, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	return LeaderId, response.Response{Status: http.StatusOK, Message: "OK"}

}

func (cr ClanRepository) GetPlayers(clanID int, txn *sql.Tx) ([]player.Player, response.Response) {
	stmt, err := txn.Prepare(`SELECT 
								p.*
							FROM 
								player p NATURAL JOIN clan_player cp
							WHERE 
								cp.clan_id = ? AND
								p.active = true;`)
	defer stmt.Close()
	if err != nil {
		return nil, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	rows, err := stmt.Query(clanID)
	if err != nil {
		return nil, response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	var players []player.Player
	var noRows bool = true
	for rows.Next() {
		noRows = false
		var player player.Player
		err := rows.Scan(
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
	if noRows {
		return nil, response.Response{Status: http.StatusNotFound, Message: "Clan Not Found"}
	}

	return players, response.Response{Status: http.StatusOK, Message: "OK"}
}
