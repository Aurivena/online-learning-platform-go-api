package pkg

import (
	"log/slog"
	"net/http"
	"online-learning-platform-go-api/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunServer(cfg config.Server, handler http.Handler) {
	httpServer := &http.Server{
		Addr:           cfg.Addr + ":" + cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
	}
	slog.Info("server started successfully")
	httpServer.ListenAndServe()
}

func StopServer() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("server is shutting down...")
}
