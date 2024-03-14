package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router *chi.Mux
	Pool *pgxpool.Pool
}

func NewApp() *App {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err.Error())
	}
	app := App{
		Router: chi.NewRouter(),
		Pool: pool,
	}

	return &app
}

func (a *App) Serve() {
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port:", port)
	http.ListenAndServe(":" + port, a.Router)
}