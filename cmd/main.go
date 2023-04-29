package main

import (
	connection "Lila-Back/Helpers/Connection"
	"Lila-Back/internal/infraestructure/playerRepository"
)

func main() {
	// TODO: get port
	txn, err := connection.Connect()
	if err != nil {
		print(err.Error())
	} else {
		player, _ := playerRepository.GetPlayer(1, txn)
		player.ToString()
	}

}
