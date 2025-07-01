package http

import (
	"github.com/gofiber/fiber"
	"github.com/tehrelt/wm-test/internal/transport/http/handlers"
	"github.com/tehrelt/wm-test/internal/usecase"
)

type Server struct {
	router *fiber.App
	uc     *usecase.UseCase
}

func New(uc *usecase.UseCase) *Server {
	return &Server{
		uc:     uc,
		router: fiber.New(),
	}
}

func (a *Server) setup() {
	a.router.Get("/", handlers.ListTasks(a.uc))
	a.router.Post("/", handlers.CreateTask(a.uc))
	a.router.Get("/:id", handlers.GetTask(a.uc))
}

func (a *Server) Run() error {
	port := 8080
	a.setup()

	return a.router.Listen(port)
}

func (a *Server) Shutdown() error {
	return a.router.Shutdown()
}
