package service

import (
	"context"
	"crypto-server/internal/client"
	"crypto-server/internal/model"
	"crypto-server/internal/utils"
	"fmt"

	"github.com/gookit/slog"
)

type ISymbol interface {
	GetAllSymbols(ctx context.Context) []string
	AddSymbol(ctx context.Context, symbol *model.Symbol) error
}

type symbol struct {
	cache        utils.ICache
	hitbtcClient client.IHitBtc
}

func NewSymbol(symbolCache utils.ICache, hitbtcClient client.IHitBtc) ISymbol {
	s := &symbol{
		cache:        symbolCache,
		hitbtcClient: hitbtcClient,
	}
	slog.Info("adding default symbols")
	// s.AddSymbol(context.TODO(), &model.Symbol{
	// 	BaseCurrency: "BTC",
	// 	PairCurrency: "USD",
	// 	Symbol:       "BTCUSD",
	// })
	s.AddSymbol(context.TODO(), &model.Symbol{
		BaseCurrency: "ETH",
		PairCurrency: "BTC",
		Symbol:       "ETHBTC",
	})
	return s
}

func (s *symbol) GetAllSymbols(ctx context.Context) []string {
	return s.cache.ReadAllKeys()
}

func (s *symbol) AddSymbol(ctx context.Context, symbol *model.Symbol) error {
	// check if it exist
	sym := s.cache.ReadVal(symbol.Symbol)
	if sym != nil {
		return fmt.Errorf("symbol %s already exist", symbol.Symbol)
	}

	exist, err := s.hitbtcClient.SymbolExist(symbol.Symbol)
	if err != nil || !exist {
		return fmt.Errorf("error fetching symbol details, error:%s", err.Error())
	}
	currency, err := s.hitbtcClient.GetCurrency(symbol.BaseCurrency)
	if err != nil {
		return fmt.Errorf("error fetching symbol details, error:%s", err.Error())
	}

	currency.ID = symbol.BaseCurrency
	currency.FeeCurrency = symbol.PairCurrency

	s.cache.Write(symbol.Symbol, *currency)

	tickerFeed, err := s.hitbtcClient.Sync(symbol.Symbol)
	if err != nil {
		slog.Errorf("failed subcribing to ticker for %s, error:%w", symbol.Symbol, err)
	} else {
		go s.updateCurrency(tickerFeed, *symbol, *currency)
	}

	return nil
}

func (s *symbol) updateCurrency(tickerFeed <-chan client.WSNotificationTickerResponse, symbol model.Symbol, currency model.Currency) {
	for {
		ticker := <-tickerFeed
		currency.Ask = ticker.Ask
		currency.Bid = ticker.Bid
		currency.Last = ticker.Last
		currency.Open = ticker.Open
		currency.Low = ticker.Low
		currency.High = ticker.High
		s.cache.Write(symbol.Symbol, currency)
	}
}
