package player

import "fmt"

type Player struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Rank    int    `json:"rank"`
	Winrate int    `json:"winrate"`
}

//TODO scar esto
func (p Player) ToString() {
	fmt.Println("ID", p.Id, "N", p.Name, p.Level, p.Rank, p.Winrate)
}
