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
		consumedNumber: 0,
	}
}

type ComputerVWAP struct {
	trades         map[string][]*Trade
	sumVolume      map[string]float64
	priceTimesSize map[string]float64
	vWAP           map[string]float64
	windowSize     int
	consumedNumber int
}

func (v *ComputerVWAP) Listen(ctx context.Context, cancelFunc context.CancelFunc, wsf WebSocketFeed) {

	// Subscribe to socket
	if err := wsf.Subscribe(); err != nil {
		return
	}
	defer wsf.TurnOff()

	dst := make(chan *Trade)
	defer close(dst)

	go func() {
		//Consume the socket
		trade, err := wsf.Read()
		if err != nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				dst <- nil
				return
			case dst <- trade:
				//Consume the socket
				trade, err = wsf.Read()
				if err != nil {
					return
				}
			}
		}
	}()

	count := 0

	for trade := range dst {
		if trade == nil {
			return
		}

		if count == 10 {
			cancelFunc()
		}
		count++

		// single compute example
		v.Compute(trade)

		log.Printf(
			"Symbol: %s Trade Sum:%3d VWAP: %s %.2f\n",
			trade.TradeSymbol, len(v.trades[trade.TradeSymbol]), trade.Currency, v.vWAP[trade.TradeSymbol],
		)
	}
}

// Compute does the main calculation formula of VWAP price.
func (v *ComputerVWAP) Compute(trade *Trade) {
	tradeSymbol := trade.TradeSymbol
	volume, _ := strconv.ParseFloat(trade.Volume, 64)
	price, _ := strconv.ParseFloat(trade.Price, 64)
	v.consumedNumber++

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
