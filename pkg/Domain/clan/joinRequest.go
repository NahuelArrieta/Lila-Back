package clan

type JoinRequest struct {
	Player_Id int    `json:"playerId"`
	Clan_Id   int    `json:"clanId"`
	Password  string `json:"password"`
}
