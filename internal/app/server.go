package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shinichi-sunayama/urlshort/internal/logic"
)

type Server struct {
	echo  *echo.Echo
	logic *logic.Service
}

func New(e *echo.Echo, logic *logic.Service) *Server {
	return &Server{
		echo:  e,
		logic: logic,
	}
}

func (s *Server) Routes() {
	// 静的ファイルとテンプレート
	s.echo.Static("/static", "templates/static")
	s.echo.Renderer = NewTemplateRenderer("templates")

	// ルーティング
	s.echo.GET("/", s.handleIndex)           // フォーム表示
	s.echo.POST("/", s.handleShorten)        // フォームPOST処理
	s.echo.GET("/r/:code", s.handleRedirect) // リダイレクト処理
}

// GET /
func (s *Server) handleIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]any{
		"ShortURL": "",
		"Origin":   "",
		"Error":    "",
	})
}

// POST /
func (s *Server) handleShorten(c echo.Context) error {
	url := c.FormValue("url")
	custom := c.FormValue("custom")

	code, err := s.logic.Shorten(c.Request().Context(), url, custom)
	if err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]any{
			"ShortURL": "",
			"Origin":   url,
			"Error":    err.Error(),
		})
	}

	shortURL := c.Scheme() + "://" + c.Request().Host + "/r/" + code
	return c.Render(http.StatusOK, "index.html", map[string]any{
		"ShortURL": shortURL,
		"Origin":   url,
		"Error":    "",
	})
}

// GET /r/:code
func (s *Server) handleRedirect(c echo.Context) error {
	code := c.Param("code")
	url, _, err := s.logic.Resolve(c.Request().Context(), code)
	if err != nil || url == "" {
		return c.String(http.StatusNotFound, "404 not found")
	}
	return c.Redirect(http.StatusMovedPermanently, url)
}
