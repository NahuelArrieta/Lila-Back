package clanRepository

import (
	"Lila-Back/pkg/Domain/clan"
	"Lila-Back/pkg/Domain/player"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) int
	JoinClan(playerID, clanID int, txn *sql.Tx) int
	GetClanPassword(clanID int, txn *sql.Tx) (string, int)
	PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) int
	GetClanLeader(clanID int, txn *sql.Tx) (int, int)
	GetPlayers(clanID int, txn *sql.Tx) ([]player.Player, int)
}

type ClanRepository struct{}

func (cr ClanRepository) CreateClan(clan *clan.Clan, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`INSERT INTO clan (
								name,
								hashed_password,
								leader_id)
							VALUES (?,?,?);`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(
		clan.Name,
		clan.HashedPassword,
		clan.LeaderId,
	)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "clan.name") {
			// TODO check if name is empty or duplicate
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}

	id, err := res.LastInsertId()
	if err != nil {
		return http.StatusInternalServerError
	}
	clan.Id = int(id)

	return http.StatusOK
}

func (cr ClanRepository) JoinClan(playerID int, clanID int, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`INSERT INTO clan_player(
								player_id,
								clan_id)
							VALUES (?,?);`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	_, err = stmt.Exec(playerID, clanID)
	if err != nil {
		str := err.Error()
		if strings.Contains(str, "Duplicate") {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}

	return http.StatusOK

}

func (cr ClanRepository) GetClanPassword(clanID int, txn *sql.Tx) (string, int) {
	stmt, err := txn.Prepare(`SELECT 
								hashed_password
							FROM 
								clan
							WHERE 
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return "", http.StatusInternalServerError
	}

	var password string
	err = stmt.QueryRow(clanID).Scan(&password)

	if err == sql.ErrNoRows {
		return "", http.StatusNotFound
	}
	if err != nil {
		return "", http.StatusInternalServerError
	}
	return password, http.StatusOK

}

func (cr ClanRepository) PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) int {
	stmt, err := txn.Prepare(`UPDATE 
								clan
							SET 
								coleader_id = ?
							WHERE
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return http.StatusInternalServerError
	}

	res, err := stmt.Exec(coleaderReq.ColeaderId, coleaderReq.ClanId)
	if err != nil {
		str := err.Error()
		// TODO change fk_constraint
		if strings.Contains(str, "foreign key") {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError
	}
	if int(rows) == 0 {
		return http.StatusNotModified
	}

	return http.StatusOK
}

func (cr ClanRepository) GetClanLeader(clanID int, txn *sql.Tx) (int, int) {
	stmt, err := txn.Prepare(`SELECT 
								leader_id
							FROM 
								clan
							WHERE 
								clan_id = ?;`)
	defer stmt.Close()
	if err != nil {
		return -1, http.StatusInternalServerError
	}

	var LeaderId int
	err = stmt.QueryRow(clanID).Scan(&LeaderId)

	if err == sql.ErrNoRows {
		return -1, http.StatusNotFound
	}
	if err != nil {
		return -1, http.StatusInternalServerError
	}
	return LeaderId, http.StatusOK

}

func (cr ClanRepository) GetPlayers(clanID int, txn *sql.Tx) ([]player.Player, int) {
	stmt, err := txn.Prepare(`SELECT 
								p.*
							FROM 
								player p NATURAL JOIN clan_player cp
							WHERE 
								cp.clan_id = ? AND
								p.active = true;`)
	defer stmt.Close()
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	rows, err := stmt.Query(clanID)
	if err != nil {
		return nil, http.StatusInternalServerError
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
			return nil, http.StatusInternalServerError
		}
		players = append(players, player)

	}
	if noRows {
		return nil, http.StatusNotFound
	}

	return players, http.StatusOK
}
