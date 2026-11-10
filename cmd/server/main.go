package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/koba-masa/equivalent_exchange_v2/internal/config"
	"github.com/koba-masa/equivalent_exchange_v2/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := config.Load()
	httpServer := &http.Server{
		Addr:    cfg.Addr(),
		Handler: server.New(),
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("サーバーを起動しました", "addr", cfg.Addr())
		errCh <- httpServer.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("サーバーが異常終了しました", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("サーバーを停止しています")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("サーバーの停止に失敗しました", "error", err)
			os.Exit(1)
		}
		logger.Info("サーバーを停止しました")
	}
}
