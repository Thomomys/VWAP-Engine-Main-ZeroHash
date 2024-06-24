package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/seedcx/vwap-engine/pkg"
)

var (
	env *EnvVar
)

func main() {
	LoadEnv()

	log.Println("starting VWAP engine")

	ws := pkg.NewWebSocket()
	coinbase := pkg.NewCoinbaseFeed(ws, env.Debug)

	v := pkg.NewComputerVWAP(env.SlidingWindow)

	ctx := context.Background()
	signalCtx, cancelFn := cancelSignal(ctx)

	v.Listen(signalCtx, cancelFn, coinbase)

	log.Println("VWAP engine stopped")
}

type EnvVar struct {
	Debug         bool
	SlidingWindow int
}

// LoadEnv environment vars
func LoadEnv() {
	e := &EnvVar{
		Debug:         false,
		SlidingWindow: 200,
	}

	d := os.Getenv("DEBUG")
	e.Debug, _ = strconv.ParseBool(d)

	w := os.Getenv("DATA_SLIDE_WINDOW")
	dsw, _ := strconv.Atoi(w)
	if dsw > 0 {
		e.SlidingWindow = dsw
	}
	env = e
}

// cancelSignal Listen OS signal to stop execution close network connection properly
func cancelSignal(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	ctxCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-osSignal
	}()
	return ctxCancel, cancel
}
