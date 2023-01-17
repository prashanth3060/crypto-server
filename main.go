package main

import (
	"context"
	"crypto-server/internal/client"
	"crypto-server/internal/handler"
	"crypto-server/internal/httpserver"
	"crypto-server/internal/service"
	"crypto-server/internal/utils"
	"os"
	"os/signal"
	"time"

	"github.com/gookit/slog"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	errGroup, ectx := errgroup.WithContext(ctx)
	cache := utils.NewCache()

	hitbtcClient, socketClose, err := client.NewHitBtc("api.hitbtc.com", "UDP", cache)
	if err != nil {
		slog.Fatal("failed to connect hitbtc, error: %w", err)
	}

	defer func() {
		socketClose()
	}()

	currency := service.NewCurrency(cache)
	symbol := service.NewSymbol(cache, hitbtcClient)
	currencyHandler := handler.NewCurrencyHandler(currency)
	symbolHandler := handler.NewSymbolHandler(symbol)

	server := httpserver.NewHTTPServer(ectx, currencyHandler.GetRoutes(), symbolHandler.GetRoutes())

	errGroup.Go(func() error {
		return server.Start()
	})

	errGroup.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Warnf("failed to shutdown http server, error: %w", err)
		}
		return ctx.Err()
	})
	err = errGroup.Wait()

	slog.Errorf("stopping http server, error: %w", err)
	err = server.Close()
	if err != nil {
		slog.Error("error closing the server, error: %w", err)
	}
}
