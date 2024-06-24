package pkg

import (
	"context"
	"log"
)

func NewComputerVWAP(windowSize int) *ComputerVWAP {
	return &ComputerVWAP{
		trades:         map[string][]*Trade{},
		sumVolume:      map[string]float64{},
		priceTimesSize: map[string]float64{},
		vWAP:           map[string]float64{},
		windowSize:     windowSize,
	}
}

type ComputerVWAP struct {
	trades         map[string][]*Trade
	sumVolume      map[string]float64
	priceTimesSize map[string]float64
	vWAP           map[string]float64
	windowSize     int
}

func (v *ComputerVWAP) Listen(_ctx context.Context, _cancelFunc context.CancelFunc, wsf WebSocketFeed) {

	// Subscribe to socket
	// TODO Subscribe and listen to consume the socket feed

	//Consume the socket
	trade, err := wsf.Read()
	if err != nil {
		log.Fatal(err)
	}

	// single compute example
	v.Compute(trade)

	log.Printf(
		"Symbol: %s Trade Sum:%3d VWAP: %s %.2f\n",
		trade.TradeSymbol, len(v.trades[trade.TradeSymbol]), trade.Currency, v.vWAP[trade.TradeSymbol],
	)

}

// Compute does the main calculation formula of VWAP price.
func (v *ComputerVWAP) Compute(trade *Trade) {
	panic("implement me")
}
