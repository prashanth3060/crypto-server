package client

import (
	"encoding/json"
	"fmt"

	"github.com/gookit/slog"
	"github.com/gorilla/websocket"
)

// responseChannels handles all incoming data from the hitbtc connection.
type responseChannels struct {
	notifications notificationChannels
	ErrorFeed     chan error
}

// notificationChannels contains all the notifications from hitbtc for subscribed feeds.
type notificationChannels struct {
	TickerFeed map[string]chan WSNotificationTickerResponse
}

// Process process all incoming connections and fills the channels properly.
func (c *WSClient) Process() {

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			slog.Error(err.Error())
		}
		resp := WsResponse{}
		err = json.Unmarshal(p, &resp)
		if err != nil {
			slog.Error("error parsing socket data")
			continue
		}

		if resp.Ch != "" {
			message := resp.Data
			switch resp.Ch {
			case "ticker/1s":
				for k, v := range message {
					c.updates.notifications.TickerFeed[k] <- v
				}
			}
		}

	}
}

// NewWSClient creates a new WSClient
func NewWSClient() (*WSClient, error) {
	slog.Debug("dialing hitbtc socket")
	conn, _, err := websocket.DefaultDialer.Dial(wsAPIURLV3, nil)
	if err != nil {
		return nil, err
	}

	handler := responseChannels{
		notifications: notificationChannels{
			TickerFeed: make(map[string]chan WSNotificationTickerResponse),
		},

		ErrorFeed: make(chan error),
	}

	return &WSClient{
		conn:    conn,
		updates: &handler,
	}, nil
}

// Close closes the Websocket connected to the hitbtc api.
func (c *WSClient) Close() {
	slog.Debug("closing the socket")
	c.conn.Close()

	for _, channel := range c.updates.notifications.TickerFeed {
		close(channel)
	}

	close(c.updates.ErrorFeed)

	c.updates.notifications.TickerFeed = make(map[string]chan WSNotificationTickerResponse)
	c.updates.ErrorFeed = make(chan error)
}

// SubscribeTicker subscribes to the specified market ticker notifications.
func (c *WSClient) SubscribeTicker(symbol string) (<-chan WSNotificationTickerResponse, error) {

	req := WsRequest{
		Method: "subscribe",
		Ch:     "ticker/1s",
		Id:     1,
		Params: WSSubscriptionRequest{
			Symbols: []string{symbol},
		},
	}
	err := c.conn.WriteJSON(req)
	if err != nil {
		return nil, fmt.Errorf("failed to subcribe hitbtc ticker, error:%w", err)
	}

	if c.updates.notifications.TickerFeed[symbol] == nil {
		c.updates.notifications.TickerFeed[symbol] = make(chan WSNotificationTickerResponse)
	}

	return c.updates.notifications.TickerFeed[symbol], nil
}
