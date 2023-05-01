package clanRepository

import (
	"Lila-Back/pkg/Domain/clan"
	"database/sql"
	"net/http"
	"strings"
)

type Repository interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) int
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
