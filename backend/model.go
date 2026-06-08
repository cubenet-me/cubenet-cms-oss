package main

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Server struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Version     string    `json:"version"`
	Status      string    `json:"status"`
	TPS         float64   `json:"tps"`
	Players     int       `json:"players"`
	MaxPlayers  int       `json:"max_players"`
	Mods        []string  `json:"mods"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Build struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	ModLoader string    `json:"mod_loader"`
	MCVersion string    `json:"mc_version"`
	ServerID  string    `json:"server_id"`
	FileHash  string    `json:"file_hash"`
	FileSize  int64     `json:"file_size"`
	Changelog string    `json:"changelog"`
	CreatedAt time.Time `json:"created_at"`
}

type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}
