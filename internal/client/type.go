package client

import (
	"crypto-server/internal/model"

	"github.com/gorilla/websocket"
)

type CurrencyResponse struct {
	FullName string `json:"full_name"`
}

func (c CurrencyResponse) toCurrency() *model.Currency {
	return &model.Currency{
		FullName: c.FullName,
	}
}

type SymbolResponse struct {
	Type               string `json:"type"`
	BaseCurrency       string `json:"base_currency"`
	QuoteCurrency      string `json:"quote_currency"`
	Status             string `json:"status"`
	QuantityIncrement  string `json:"quantity_increment"`
	TickSize           string `json:"tick_size"`
	TakeRate           string `json:"take_rate"`
	MakeRate           string `json:"make_rate"`
	FeeCurrency        string `json:"fee_currency"`
	MarginTrading      bool   `json:"margin_trading"`
	MaxInitialLeverage string `json:"max_initial_leverage"`
}

type WsAck struct {
	Result Result `json:"result"`
	ID     int    `json:"id"`
}

type Result struct {
	Ch            string   `json:"ch"`
	Subscriptions []string `json:"subscriptions"`
}

type WSNotificationTickerResponse struct {
	Ask       string `json:"a"`
	Bid       string `json:"b"`
	Last      string `json:"c"`
	Open      string `json:"o"`
	Low       string `json:"l"`
	High      string `json:"h"`
	LastTrade uint   `json:"L"`
}
type WsResponse struct {
	Ch   string                                  `json:"ch"`
	Data map[string]WSNotificationTickerResponse `json:"data"`
}

type WsRequest struct {
	Method string                `json:"method"`
	Ch     string                `json:"ch"`
	Params WSSubscriptionRequest `json:"params"`
	Id     int                   `json:"id"`
}

// WSSubscriptionRequest is request type on websocket subscription.
type WSSubscriptionRequest struct {
	Symbols []string `json:"symbols"`
}

type WSClient struct {
	conn    *websocket.Conn
	updates *responseChannels
}
