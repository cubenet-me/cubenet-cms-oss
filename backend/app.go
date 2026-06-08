package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool   *pgxpool.Pool
	s3     *Storage
	hub    *Hub
	secret string
	wsCfg  WSConfig
}
