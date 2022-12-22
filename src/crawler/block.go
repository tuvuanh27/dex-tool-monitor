package crawler

import (
	"context"
	"dex-tool/src/configs/db"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type CrawlBlock struct {
	Chain  string
	Client *ethclient.Client
	Ctx    context.Context
}

type RedisBlock struct {
	Number    uint64 `json:"number" bson:"number"`
	Timestamp uint64 `json:"timestamp" bson:"timestamp"`
}

func NewCrawlBlock(chain string, client *ethclient.Client, ctx context.Context) *CrawlBlock {
	return &CrawlBlock{
		Chain:  chain,
		Client: client,
		Ctx:    ctx,
	}
}

func GetKeyBlock(chain string) string {
	return fmt.Sprintf("BLOCK_NUMBER_%s", chain)
}

func (c *CrawlBlock) GetBlockNumber() {
	var redis = db.GetRedisInstance()
	redisKey := GetKeyBlock(c.Chain)

	headers := make(chan *types.Header)
	sub, _ := c.Client.SubscribeNewHead(c.Ctx, headers)
	for {
		select {
		case err := <-sub.Err():
			panic(err)
		case header := <-headers:
			block, err := c.Client.BlockByHash(c.Ctx, header.Hash())
			if err != nil {
				panic(err)
			}
			redisBlock := RedisBlock{
				Number:    block.Number().Uint64(),
				Timestamp: block.Time(),
			}
			value, err := json.Marshal(&redisBlock)
			if err != nil {
				panic(err)
			}

			log.Println("set redis key", redisKey, "value", string(value))

			err = redis.Set(redisKey, string(value), 0)
			if err != nil {
				panic(err)
			}
		}
	}
}
