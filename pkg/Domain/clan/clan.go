package clan

type Clan struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	HashedPassword string `json:"password"`
	LeaderId       int    `json:"leaderId"`
	ColeaderId     int    `json:"coleaderId"`
}
