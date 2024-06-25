package pkg

import (
	"context"
	"log"
	"strconv"
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
	if err := wsf.Subscribe(); err != nil {
		return
	}

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
	tradeSymbol := trade.TradeSymbol
	volume, _ := strconv.ParseFloat(trade.Volume, 64)
	price, _ := strconv.ParseFloat(trade.Price, 64)

	if len(v.trades[tradeSymbol]) >= v.windowSize {
		firstVolume, _ := strconv.ParseFloat(v.trades[tradeSymbol][0].Volume, 64)
		firstPrice, _ := strconv.ParseFloat(v.trades[tradeSymbol][0].Price, 64)
		v.sumVolume[tradeSymbol] -= firstVolume
		v.priceTimesSize[tradeSymbol] -= firstPrice * firstVolume
		v.trades[tradeSymbol] = v.trades[tradeSymbol][1:]
	}
	v.sumVolume[tradeSymbol] += volume
	v.priceTimesSize[tradeSymbol] += volume * price
	v.vWAP[tradeSymbol] = v.priceTimesSize[tradeSymbol] / v.sumVolume[tradeSymbol]
	v.trades[tradeSymbol] = append(v.trades[tradeSymbol], trade)
}
