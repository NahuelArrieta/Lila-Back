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
	JoinClan(jr *clan.JoinRequest, txn *sql.Tx) int
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

	return ch.Repository.JoinClan(clan.LeaderId, clan.Id, txn)
}

func (ch ClanHandler) JoinClan(jr *clan.JoinRequest, txn *sql.Tx) int {

	hashedPassword, status := ch.Repository.GetClanPassword(jr.Clan_Id, txn)
	if status != http.StatusOK {
		return status
	}

	if !hashPass.CheckPassword(jr.Password, hashedPassword) {
		return http.StatusForbidden
	}
	jr.Password = hashedPassword

	return ch.Repository.JoinClan(jr.Player_Id, jr.Clan_Id, txn)
}
