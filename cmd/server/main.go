package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/shinichi-sunayama/urlshort/internal/app"
	"github.com/shinichi-sunayama/urlshort/internal/logic"
	"github.com/shinichi-sunayama/urlshort/internal/store"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_ = godotenv.Load() // .env が無くても起動可

	// 依存の組み立て
	st, err := store.NewRedisStore()
	must(err)
	svc := &logic.Service{Store: st}

	e := echo.New()
	s := app.New(e, svc)
	s.Routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	must(e.Start(":" + port))
}
