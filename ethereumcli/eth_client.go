package ethereumcli

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func EthClientWithTimeout(ctx context.Context, url string) (*ethclient.Client, error) {
	ctxt, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return ethclient.DialContext(ctxt, url)
}
