package clanRepository

import (
	"Lila-Back/pkg/Domain/clan"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) int
	JoinClan(playerID int, clanID int, txn *sql.Tx) int
	GetClanPassword(clanID int, txn *sql.Tx) (string, int)
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
							VALUES (?,?)`)
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
