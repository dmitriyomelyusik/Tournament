// Package entity contains all entities, that used in application
package entity

// Player is struct for players perfomance
type Player struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
}

// Winner is player, who won tournament
type Winner struct {
	ID     string `json:"id"`
	Points int    `json:"points"`
	Prize  int    `json:"prize"`
}

// Winners contains every winner from tournaments
type Winners struct {
	Winners []Winner `json:"winners"`
}

// Tournament is struct for tournament perfomance
type Tournament struct {
	ID           string   `json:"id"`
	Deposit      int      `json:"deposit"`
	Prize        int      `json:"prize"`
	Participants []string `json:"participants"`
	Winner       Winner   `json:"winner"`
	IsOpen       bool     `json:"isOpen"`
}
