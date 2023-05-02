package player

type Player struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Level   int     `json:"level"`
	Rank    int     `json:"rank"`
	Winrate float32 `json:"winrate"`
	Active  bool    `json:"active"`
}

func (player Player) SetDefaultValues() {
	player.Level = 0
	player.Rank = 0
	player.Winrate = 0
	player.Active = true
}
