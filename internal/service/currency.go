package service

import (
	"context"
	"crypto-server/internal/model"
	"crypto-server/internal/utils"
	"fmt"
)

type ICurrency interface {
	GetAllCurrency(ctx context.Context) (*[]model.Currency, error)
	GetCurrencyBySymbol(ctx context.Context, symbol string) (*model.Currency, error)
}

type currency struct {
	cache utils.ICache
}

func NewCurrency(symbolCache utils.ICache) ICurrency {
	return &currency{
		cache: symbolCache,
	}
}

func (c *currency) GetAllCurrency(ctx context.Context) (*[]model.Currency, error) {
	currencies := c.cache.ReadVal("")
	if currencies != nil {
		data, ok := currencies.([]model.Currency)
		if !ok {
			return nil, fmt.Errorf("invalid data")
		}
		return &data, nil
	}
	// go call live data
	return nil, nil
}

func (c *currency) GetCurrencyBySymbol(ctx context.Context, symbol string) (*model.Currency, error) {
	currencies := c.cache.ReadVal(symbol)
	if currencies != nil {
		data, ok := currencies.(model.Currency)
		if !ok {
			return nil, fmt.Errorf("invalid data")
		}
		return &data, nil
	}
	// go call live data
	return nil, nil
}
