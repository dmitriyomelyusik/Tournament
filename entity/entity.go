// Package entity contains all entities, that used in application
package entity

// Player is struct for players perfomance
type Player struct {
	ID     string `json:"id" bson:"_id"`
	Points int    `json:"points" bson:"points"`
}

// Winner is player, who won tournament
type Winner struct {
	ID     string `json:"id" bson:"_id"`
	Points int    `json:"points" bson:"points"`
	Prize  int    `json:"prize" bson:"prize"`
}

// Winners contains every winner from tournaments
type Winners struct {
	Winners []Winner `json:"winners" bson:"winners"`
}

// Tournament is struct for tournament perfomance
type Tournament struct {
	ID           string   `json:"id" bson:"_id"`
	Deposit      int      `json:"deposit" bson:"deposit"`
	Prize        int      `json:"prize" bson:"prize"`
	Participants []string `json:"participants" bson:"participants"`
	Winner       Winner   `json:"winner" bson:"winner"`
	IsOpen       bool     `json:"isOpen" bson:"isOpen"`
}
