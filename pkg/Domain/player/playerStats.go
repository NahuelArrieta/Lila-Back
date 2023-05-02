package player

type PlayerStats struct {
	Level   int     `json:"level,required"`
	Rank    int     `json:"rank,required"`
	Winrate float32 `json:"winrate,required"`
}

// func (ps PlayerStats) Init() PlayerStats {
// 	ps.Level = -1
// 	ps.Rank = -1
// 	ps.Winrate = -1
// 	return ps
// }
