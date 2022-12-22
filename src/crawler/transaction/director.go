package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Director struct {
	builder ITransaction
}

func (d *Director) SetBuilder(builder ITransaction) {
	d.builder = builder
}

func (d *Director) BuildTransaction(chain string, client *ethclient.Client, ctx context.Context, pairAbi string, topic string) Transaction {
	d.builder.SetChain(chain).SetClient(client).SetCtx(ctx).SetPairAbi(pairAbi).SetTopic(topic)
	return d.builder.GetTransaction()
}
