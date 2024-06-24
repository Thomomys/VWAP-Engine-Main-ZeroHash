package pkg

import (
	"encoding/json"
	"errors"

	gws "github.com/gorilla/websocket"
)

var (
	ErrStopped = errors.New("websocket service stopped")
)

type WebSocket struct {
	IsConnected bool
	NumberAsStr bool
	Conn        *gws.Conn
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		IsConnected: false,
		NumberAsStr: true,
	}
}

func (c *WebSocket) WriteJSON(message interface{}) error {
	if !c.IsConnected {
		return ErrStopped
	}
	if err := c.Conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (c *WebSocket) ReadJSON(message interface{}) error {
	if !c.IsConnected {
		return ErrStopped
	}
	_, r, err := c.Conn.NextReader()
	if err != nil {
		return err
	}
	d := json.NewDecoder(r)

	if c.NumberAsStr {
		d.UseNumber()
	}
	return d.Decode(message)
}

func (c *WebSocket) Connect(endpoint string) error {
	if c.Conn != nil {
		c.IsConnected = true
		return nil
	}

	conn, _, err := gws.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return err
	}

	c.Conn = conn
	c.IsConnected = true
	return nil
}

func (c *WebSocket) Disconnect() error {
	c.IsConnected = false
	if err := c.Conn.Close(); err != nil {
		return err
	}
	return nil
}
