package http

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers"
	"github.com/tehrelt/wm-test/internal/usecase"
)

type Server struct {
	router *echo.Echo
	uc     *usecase.UseCase
}

func New(uc *usecase.UseCase) *Server {
	return &Server{
		uc:     uc,
		router: echo.New(),
	}
}

func (a *Server) setup() {
	a.router.GET("/", handlers.ListTasks(a.uc))
	a.router.POST("/", handlers.CreateTask(a.uc))
	a.router.GET("/:id", handlers.GetTask(a.uc))
	a.router.DELETE("/:id", handlers.DeleteTask(a.uc))
}

func (a *Server) Run() error {
	port := 8080
	a.setup()

	return a.router.Start(fmt.Sprintf(":%d", port))
}

func (a *Server) Shutdown(ctx context.Context) error {
	return a.router.Shutdown(ctx)
}
