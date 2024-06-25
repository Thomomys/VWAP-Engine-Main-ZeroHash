package pkg

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewCoinbaseFeed(t *testing.T) {
	type args struct {
		ws    *WebSocket
		debug bool
	}
	tests := []struct {
		name string
		args args
		want WebSocketFeed
	}{
		{
			name: "",
			args: args{
				ws:    NewWebSocket(),
				debug: false,
			},
			want: &CoinbaseFeed{
				ws:    NewWebSocket(),
				debug: false,
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCoinbaseFeed(tt.args.ws, tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCoinbaseFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinbaseFeed_Read(t *testing.T) {
	type fields struct {
		Conf *CoinbaseConfig
		ws   *WebSocket
	}
	tests := []struct {
		name        string
		fields      fields
		want        *Trade
		wantErr     bool
		wantErrType error
		turnOff     bool
	}{
		{
			name: "fail on read with not connected websocket operation",
			fields: fields{
				Conf: &CoinbaseConfig{},
				ws: &WebSocket{
					IsConnected: false,
					NumberAsStr: true,
				},
			},
			want:        nil,
			wantErr:     true,
			wantErrType: ErrStopped,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CoinbaseFeed{
				conf: tt.fields.Conf,
				ws:   tt.fields.ws,
			}
			got, err := c.Read()
			if tt.wantErr && !errors.Is(err, tt.wantErrType) {
				t.Errorf("CoinbaseFeed.Read() error =%+v, wantErr (%+v)", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoinbaseFeed.Read() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
