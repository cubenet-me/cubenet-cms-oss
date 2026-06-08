package model

type Server struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Version     string   `json:"version"`
	Status      string   `json:"status"`
	TPS         float64  `json:"tps"`
	Players     int      `json:"players"`
	MaxPlayers  int      `json:"max_players"`
	Mods        []string `json:"mods"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
