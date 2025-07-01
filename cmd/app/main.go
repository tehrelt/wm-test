package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tehrelt/wm-test/internal/app"
	"github.com/tehrelt/wm-test/internal/storage/memo"
	"github.com/tehrelt/wm-test/internal/transport/http"
	"github.com/tehrelt/wm-test/internal/usecase"
)

func main() {
	st := memo.New()
	uc := usecase.New(st)
	server := http.New(uc)
	app := app.New(server)

	start := time.Now()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-sigchan

	slog.Info("server stopped", slog.Duration("duration", time.Since(start)), slog.String("signal", sig.String()))
}
