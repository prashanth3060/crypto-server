package client

import (
	"crypto-server/internal/model"
	"crypto-server/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gookit/slog"
)

const (
	wsAPIURLV3 string = "wss://api.hitbtc.com/api/3/ws/public"
	httpAPIURL string = "https://api.hitbtc.com/api/3/public"
)

type IHitBtc interface {
	// Dial() (net.Conn, error)
	Sync(symbol string) (<-chan WSNotificationTickerResponse, error)
	GetCurrency(currency string) (*model.Currency, error)
	SymbolExist(symbol string) (bool, error)
}

type hitbtc struct {
	wsconn *WSClient
	cache  utils.ICache
}

func NewHitBtc(host, connectionType string, cache utils.ICache) (IHitBtc, func(), error) {
	wsc, err := NewWSClient()
	go wsc.Process()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to %s, error: %s", host, err.Error())
	}
	fn := wsc.Close
	return &hitbtc{
		wsconn: wsc,
		cache:  cache,
	}, fn, nil
}

func (h *hitbtc) Sync(symbol string) (<-chan WSNotificationTickerResponse, error) {

	tickerFeed, err := h.wsconn.SubscribeTicker(symbol)
	if err != nil {
		return nil, err
	}

	return tickerFeed, nil
}

func (h *hitbtc) SymbolExist(symbol string) (bool, error) {
	slog.Debugf("getting symbol details from hitbtc, symbol: %s", symbol)

	url := httpAPIURL + "/symbol?symbols=" + symbol

	req, err := utils.CreateRequest(http.MethodGet, url, "", nil)
	if err != nil {
		slog.Errorf("failed to create get currency request, error:%w", err)
		return false, fmt.Errorf("failed to get symbol details")
	}

	resp, err := utils.SendRequest(req, http.StatusOK)
	if err != nil {
		slog.Errorf("failed to send get currency request, error:%w", err)
		return false, fmt.Errorf("failed to get symbol details")
	}
	c := make(map[string]SymbolResponse)
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		slog.Errorf("error parsing get currency response, error:%w", err)
		return false, fmt.Errorf("failed to get symbol details")
	}

	_, ok := c[symbol]
	if !ok {
		return false, fmt.Errorf("symbol not found")
	}
	slog.Debugf("found currency details from hitbtc, currency: %s", symbol)
	return true, nil
}

func (h *hitbtc) GetCurrency(currency string) (*model.Currency, error) {

	slog.Debugf("getting currency details from hitbtc, currency: %s", currency)

	url := httpAPIURL + "/currency?currencies=" + currency

	req, err := utils.CreateRequest(http.MethodGet, url, "", nil)
	if err != nil {
		slog.Errorf("failed to create get currency request, error:%w", err)
		return nil, fmt.Errorf("failed to get symbol details")
	}

	resp, err := utils.SendRequest(req, http.StatusOK)
	if err != nil {
		slog.Errorf("failed to send get currency request, error:%w", err)
		return nil, fmt.Errorf("failed to get symbol details")
	}
	c := make(map[string]CurrencyResponse)
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		slog.Errorf("error parsing get currency response, error:%w", err)
		return nil, fmt.Errorf("failed to get symbol details")
	}

	v, ok := c[currency]
	if !ok {
		return nil, fmt.Errorf("symbol not found")
	}
	slog.Debugf("found currency details from hitbtc, currency: %s", currency)
	return v.toCurrency(), nil
}
