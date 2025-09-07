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

func main() {
	// .env ロード（.env はこの main.go と同じ cmd/server/ に置く）
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠ .env ファイルが読み込めませんでした: %v", err)
	}

	// Redis ストア初期化
	st, err := store.NewRedisStore()
	if err != nil {
		log.Fatal(err)
	}

	// サービス層（ビジネスロジック）
	svc := &logic.Service{Store: st}

	// Echo インスタンス作成
	e := echo.New()

	// サーバー構造体にサービスと Echo を渡す
	s := app.New(e, svc)
	s.Routes()

	// ポート指定して起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 サーバー起動: http://localhost:%s", port)
	log.Fatal(e.Start(":" + port))
}
