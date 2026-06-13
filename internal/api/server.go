package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fareez-ahamed/go-ledger-rest/internal/config"
)

func Run(ctx context.Context, cfg *config.Config) error {
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: NewRouter(),
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}
}
