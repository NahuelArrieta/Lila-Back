package clanHandler

import (
	"Lila-Back/internal/infraestructure/clanRepository"
	"Lila-Back/pkg/Domain/clan"
	"Lila-Back/pkg/Domain/response"
	hashPass "Lila-Back/pkg/Helpers/HashPass"
	"database/sql"
	"net/http"
)

type Handler interface {
	CreateClan(clan *clan.Clan, txn *sql.Tx) response.Response
	JoinClan(jr *clan.JoinRequest, txn *sql.Tx) response.Response
	PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) response.Response
	GetPlayers(clanID int, txn *sql.Tx) (interface{}, response.Response)
}

type ClanHandler struct {
	Repository clanRepository.Repository
}

func (ch ClanHandler) CreateClan(clan *clan.Clan, txn *sql.Tx) response.Response {
	hashedPassword, err := hashPass.HashPassword(clan.HashedPassword)
	if err != nil {
		return response.Response{Status: http.StatusInternalServerError, Message: "Internal Server Error"}
	}
	clan.HashedPassword = hashedPassword

	resp := ch.Repository.CreateClan(clan, txn)
	if resp.Status != http.StatusOK {
		return resp
	}

	resp = ch.Repository.JoinClan(clan.LeaderId, clan.Id, txn)

	return resp
}

func (ch ClanHandler) JoinClan(jr *clan.JoinRequest, txn *sql.Tx) response.Response {

	hashedPassword, resp := ch.Repository.GetClanPassword(jr.Clan_Id, txn)
	if resp.Status != http.StatusOK {
		return resp
	}

	if !hashPass.CheckPassword(jr.Password, hashedPassword) {
		return response.Response{Status: http.StatusForbidden, Message: "Password Incorrect"}
	}
	jr.Password = hashedPassword

	return ch.Repository.JoinClan(jr.Player_Id, jr.Clan_Id, txn)
}

func (ch ClanHandler) PutColeader(coleaderReq clan.ColeaderRequest, txn *sql.Tx) response.Response {
	leaderId, resp := ch.Repository.GetClanLeader(coleaderReq.ClanId, txn)
	if resp.Status != http.StatusOK {
		return resp
	}
	if leaderId != coleaderReq.LeaderId {
		return response.Response{Status: http.StatusForbidden, Message: "Leader Incorrect"}
	}

	// Join the coleader to the clan in case is is not in.
	resp = ch.Repository.JoinClan(coleaderReq.ColeaderId, coleaderReq.ClanId, txn)
	if resp.Status != http.StatusOK {
		return resp
	}

	return ch.Repository.PutColeader(coleaderReq, txn)
}

func (ch ClanHandler) GetPlayers(clanID int, txn *sql.Tx) (interface{}, response.Response) {
	return ch.Repository.GetPlayers(clanID, txn)
}
