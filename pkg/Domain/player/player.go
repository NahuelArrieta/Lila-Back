package player

type Player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Rank    int    `json:"rank"`
	Winrate int    `json:"winrate"`
	Active  bool   `json:"active"`
}
