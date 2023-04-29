package player

import "fmt"

type Player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Level   string `json:"level"`
	Rank    string `json:"rank"`
	Winrate string `json:"winrate"`
}

//TODO scar esto
func (p Player) ToString() {
	fmt.Println("ID", p.Id, "N", p.Name, p.Level, p.Rank, p.Winrate)
}
