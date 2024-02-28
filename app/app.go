package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

type App struct {
	Router *chi.Mux
	Pool *pgxpool.Pool
}

func NewApp() *App {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
	app := App{
		Router: chi.NewRouter(),
		Pool: pool,
	}

	return &app
}

func (a *App) Serve() {
	http.ListenAndServe(":" + os.Getenv("PORT"), a.Router)
}