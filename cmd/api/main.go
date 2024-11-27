package main

import (
	"log"

	"github.com/Gustavicho/gosocial/internal/db"
	"github.com/Gustavicho/gosocial/internal/env"
	"github.com/Gustavicho/gosocial/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	godotenv.Load()
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),

		db: dbConfig{
			dsn:          env.GetString("DB_DSN", "postgres://root:password@localhost:5432/gosocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection established")

	app := &application{
		config: cfg,
		store:  store.NewStorage(db),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
