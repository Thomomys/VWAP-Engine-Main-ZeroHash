package pkg

// Trade represents the main trading operation
type Trade struct {
	Id           int64
	Volume       string
	Price        string
	TradeSymbol  string
	ProviderName string
	Currency     string
}

// WebSocketFeed defines minimal websocket functions to *Trade feed
type WebSocketFeed interface {
	Subscribe() error
	TurnOff() error
	Read() (*Trade, error)
}

// WebSocketService wrapper to websocket library
type WebSocketService interface {
	Disconnect()
	Connect(endpoint string) error
	ReadJSON(message interface{}) error
	WriteJSON(message interface{}) error
}
