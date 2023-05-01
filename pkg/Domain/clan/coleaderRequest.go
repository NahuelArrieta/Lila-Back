package clan

type ColeaderRequest struct {
	ClanId     int `json:"id"`
	LeaderId   int `json:"leaderId"`
	ColeaderId int `json:"coleaderId"`
}
