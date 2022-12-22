package pair

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

type pairCrawler[T any] struct {
	Chain    string
	Client   *ethclient.Client
	Ctx      context.Context
	Instance *T
}
