package config

import (
	"log"
	"os"
)

type Config struct {
	Port     string
	DBDSN    string
	MediaDir string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "retrobytes.db"
	} // sqlite file in project root
	media := os.Getenv("MEDIA_DIR")
	if media == "" {
		media = "./media"
	}

	cfg := Config{Port: port, DBDSN: dsn, MediaDir: media}
	log.Printf("[config] PORT=%s DB_DSN=%s MEDIA_DIR=%s", cfg.Port, cfg.DBDSN, cfg.MediaDir)
	return cfg
}
