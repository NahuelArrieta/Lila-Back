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
	PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) int
	GetPlayers(clanID int, txn *sql.Tx) (interface{}, int)
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

	status = ch.Repository.JoinClan(clan.LeaderId, clan.Id, txn)

	return status
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

func (ch ClanHandler) PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) int {
	leaderId, status := ch.Repository.GetClanLeader(coleaderReq.ClanId, txn)
	if status != http.StatusOK {
		return status
	}
	if leaderId != coleaderReq.LeaderId {
		return http.StatusForbidden
	}

	// Join the coleader to the clan in case is is not in.
	status = ch.Repository.JoinClan(coleaderReq.ColeaderId, coleaderReq.ClanId, txn)
	if status != http.StatusBadRequest && status != http.StatusOK {
		// BadRequest is allowed in case the player is already in the clan
		return status
	}

	return ch.Repository.PutColeader(coleaderReq, txn)
}

func (ch ClanHandler) GetPlayers(clanID int, txn *sql.Tx) (interface{}, int) {
	return ch.Repository.GetPlayers(clanID, txn)
}
