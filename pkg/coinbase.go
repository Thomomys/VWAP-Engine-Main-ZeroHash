package pkg

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Message represents Coinbase Channel WebSocket Message
type Message struct {
	ProductID string           `json:"product_id"`
	Price     string           `json:"price"`
	Size      string           `json:"size"`
	Type      string           `json:"type"`
	TradeID   int64            `json:"trade_id"`
	Sequence  int64            `json:"sequence"`
	Side      string           `json:"side"`
	Time      time.Time        `json:"time"`
	Channels  []ChannelMessage `json:"channels"`
	Message   string           `json:"message,omitempty"`
	Reason    string           `json:"reason,omitempty"`
}

type ChannelMessage struct {
	Name        string   `json:"name"`
	ProductsIDs []string `json:"product_ids"`
}

// CoinbaseConfig represents all available settings for coinbase websocket
type CoinbaseConfig struct {
	ProviderName    string
	Endpoint        string
	Currency        string
	NumberAsStr     bool
	Pairs           []string
	SubRetries      int
	SubChannelName  string
	SubResponseType string
}

// CoinbaseFeed interacts with websocket server operations
type CoinbaseFeed struct {
	conf  *CoinbaseConfig
	ws    *WebSocket
	debug bool
}

// NewCoinbaseFeed takes a conf and internal/network.WebSockt service to
// subscribe and listen the websocket provider
func NewCoinbaseFeed(ws *WebSocket, debug bool) WebSocketFeed {
	return &CoinbaseFeed{
		ws:    ws,
		debug: debug,
		conf: &CoinbaseConfig{
			ProviderName:    "coinbase.com",
			Endpoint:        "wss://ws-feed.exchange.coinbase.com",
			Currency:        "USD",
			NumberAsStr:     true,
			Pairs:           []string{"BTC-USD", "ETH-USD", "ETH-BTC"},
			SubRetries:      30,
			SubChannelName:  "matches",
			SubResponseType: "subscriptions",
		},
	}
}

// Subscribe to 'c.conf.Pairs' on the channel 'c.conf.SubChannelName'
func (c *CoinbaseFeed) Subscribe() error {
	if err := c.ws.Connect(c.conf.Endpoint); err != nil {
		log.Fatalf("fail to connect into websocket=%s err=%s", c.conf.Endpoint, err)
		return err
	}

	sub := &Message{
		Type: "subscribe",
		Channels: []ChannelMessage{
			{
				Name:        c.conf.SubChannelName,
				ProductsIDs: c.conf.Pairs,
			},
		},
	}

	if err := c.writeAndWait(sub); err != nil {
		log.Printf("fail to subscribe channel=%s err=%s", c.conf.SubChannelName, err)
		return err
	}
	return nil
}

func (c *CoinbaseFeed) TurnOff() error {
	unSub := &Message{
		Type: "unsubscribe",
		Channels: []ChannelMessage{
			{
				Name:        c.conf.SubChannelName,
				ProductsIDs: []string{},
			},
		},
	}
	if err := c.writeAndWait(unSub); err != nil {
		log.Printf("error to unsubscribe channel=%s | err=%s", c.conf.SubChannelName, err)
		return err
	}

	if err := c.ws.Disconnect(); c.ws.IsConnected && err != nil {
		log.Printf("fail to disconnect into websocket=%s err=%s", c.conf.Endpoint, err)
		return err
	}
	return nil
}

// Read subscribed messages
func (c *CoinbaseFeed) Read() (*Trade, error) {
	msg := &Message{}
	err := c.ws.ReadJSON(&msg)

	if errors.Is(err, ErrStopped) {
		log.Println("websocket stopped")
		return nil, ErrStopped
	}

	if err != nil {
		return nil, fmt.Errorf("websocket read fail: %s", err)
	}

	return &Trade{
		TradeSymbol:  msg.ProductID,
		Volume:       msg.Size,
		Price:        msg.Price,
		ProviderName: c.conf.ProviderName,
		Currency:     c.conf.Currency,
	}, nil
}

// Interacts with the channel requesting msg.Type and msg.Channels
func (c *CoinbaseFeed) writeAndWait(msg *Message) error {
	if err := c.ws.WriteJSON(&msg); err != nil {
		return err
	}
	if err := c.waitResponse(msg); err != nil {
		return err
	}
	return nil
}

// Listen to channel messages until success response
// Fails after exceed the 'c.conf.SubRetries'
func (c *CoinbaseFeed) waitResponse(asked *Message) error {
	for i := 1; i < c.conf.SubRetries && c.ws.IsConnected; i++ {
		resp := &Message{}
		if err := c.ws.ReadJSON(&resp); err != nil {
			return fmt.Errorf("websocket read fail: %s", err)
		}

		if resp.Type == "error" {
			return fmt.Errorf("%s reason: %s", resp.Message, resp.Reason)

		}

		if resp.Type == c.conf.SubResponseType {
			if c.debug {
				log.Printf("successful %s", asked.Type)
			}
			return nil
		}
	}
	return fmt.Errorf("fail to %s:%s after %d attemps",
		asked.Type, asked.Channels[0].Name, c.conf.SubRetries,
	)
}
