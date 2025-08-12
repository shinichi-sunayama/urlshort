package app

import (
	"context"
	_ "embed" // go:embed を使うため
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ビジネスロジックの契約
type coreService interface {
	Shorten(ctx context.Context, inURL, custom string) (string, error)
	Resolve(ctx context.Context, code string) (string, bool, error)
}

type Server struct {
	E       *echo.Echo
	svc     coreService
	baseURL string
}

func New(e *echo.Echo, svc coreService) *Server {
	return &Server{
		E:       e,
		svc:     svc,
		baseURL: getenv("BASE_URL", "http://localhost:8080"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type shortenReq struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code"`
}
type shortenRes struct {
	Code   string `json:"code"`
	Short  string `json:"short_url"`
	Origin string `json:"origin_url"`
}

//go:embed docs/swagger.json
var swaggerJSON []byte

func (s *Server) Routes() {
	// ミドルウェア（ログ・リカバリ・簡易レート制限 60 req/min）
	s.E.Use(middleware.Logger())
	s.E.Use(middleware.Recover())
	s.E.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(60)))

	// 健康チェック
	s.E.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	// OpenAPI (json)
	s.E.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", swaggerJSON)
	})

	// Swagger UI (GUI) → /docs で閲覧
	s.E.GET("/docs", s.handleDocs)

	// API 本体
	s.E.POST("/shorten", s.handleShorten)
	s.E.GET("/r/:code", s.handleResolve)
	s.E.HEAD("/r/:code", s.handleResolveHEAD)
}

func (s *Server) handleShorten(c echo.Context) error {
	var req shortenReq
	if err := c.Bind(&req); err != nil || req.URL == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}
	code, err := s.svc.Shorten(c.Request().Context(), req.URL, req.CustomCode)
	if err != nil {
		// 簡易に 400 に集約（READMEで方針明記）
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, shortenRes{
		Code: code, Short: s.baseURL + "/r/" + code, Origin: req.URL,
	})
}

func (s *Server) handleResolve(c echo.Context) error {
	code := c.Param("code")
	dest, _, err := s.svc.Resolve(c.Request().Context(), code)
	if err != nil || dest == "" {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	return c.Redirect(http.StatusMovedPermanently, dest)
}

func (s *Server) handleResolveHEAD(c echo.Context) error {
	code := c.Param("code")
	dest, _, err := s.svc.Resolve(c.Request().Context(), code)
	if err != nil || dest == "" {
		return c.NoContent(http.StatusNotFound)
	}
	c.Response().Header().Set(echo.HeaderLocation, dest)
	return c.NoContent(http.StatusMovedPermanently) // 301 + Location
}

// Swagger UI のシンプル配信（CDNからJS/CSS読込、/swagger.json を表示）
func (s *Server) handleDocs(c echo.Context) error {
	html := `<!doctype html>
<html lang="ja">
<head>
<meta charset="utf-8">
<title>URL Shortener API Docs</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
<style>body { margin:0; padding:0; }</style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = () => {
  window.ui = SwaggerUIBundle({
    url: '/swagger.json',
    dom_id: '#swagger-ui',
    presets: [SwaggerUIBundle.presets.apis],
    layout: 'BaseLayout'
  });
};
</script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}
