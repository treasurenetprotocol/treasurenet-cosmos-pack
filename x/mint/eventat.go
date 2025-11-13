package mint

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func getMintat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := ethclient.Dial("http://127.0.0.1:8555") // ws://127.0.0.1:8546
			if err != nil {
			}
			defer conn.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}
