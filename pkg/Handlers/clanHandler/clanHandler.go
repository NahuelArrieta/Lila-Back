package clanHandler

import (
	"Lila-Back/internal/infraestructure/clanRepository"
	"Lila-Back/pkg/Domain/clan"
	hashPass "Lila-Back/pkg/Helpers/HashPass"
	"database/sql"
	"net/http"
)

type Handler interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) int
}

type ClanHandler struct {
	Repository clanRepository.Repository
}

func (ch ClanHandler) CreateClan(clan *clan.Clan, txn *sql.Tx) int {
	hashedPassword, err := hashPass.HashPassword(clan.HashedPassword)
	if err != nil {
		return http.StatusInternalServerError
	}
	clan.HashedPassword = hashedPassword

	status := ch.Repository.CreateClan(clan, txn)
	if status != http.StatusOK {
		return status
	}

	// TODO join leader to clan
	return status
}
