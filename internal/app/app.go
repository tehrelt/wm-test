package app

import "github.com/tehrelt/wm-test/internal/transport/http"

type App struct {
	*http.Server
}

func New(s *http.Server) *App {
	return &App{Server: s}
}
