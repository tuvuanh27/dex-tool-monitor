package transaction

import (
	"context"
	"dex-tool/src/models/entities"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITransaction interface {
	SetChain(chain string) ITransaction
	SetClient(client *ethclient.Client) ITransaction
	SetCtx(ctx context.Context) ITransaction
	SetPairAbi(pairAbi string) ITransaction
	SetTopic(topic string) ITransaction
	GetTransaction() Transaction

	CrawlTx(fromBlock, toBlock int64) (<-chan TxCrawled, int)
	HandleTx(txChanel <-chan TxCrawled, results chan<- *entities.Transaction)
}
