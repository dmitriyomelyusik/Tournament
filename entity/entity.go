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
