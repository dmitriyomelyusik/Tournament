// Package entity contains all entities, that used in application
package entity

// Player is struct for players perfomance
type Player struct {
	ID     string `json:"id"`
	Points uint   `json:"points"`
}

// Winner is player, who won tournament
type Winner struct {
	ID     string `json:"winner"`
	Points uint   `json:"balance"`
	Prize  uint   `json:"prize"`
}
