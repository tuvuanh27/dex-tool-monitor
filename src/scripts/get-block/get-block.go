package main

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/configs/env"
	"dex-tool/src/crawler"
	"dex-tool/src/models/repositories"
	_chain "dex-tool/src/utils/chain"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
	"strings"
)

func main() {
	argsWithProg := os.Args
	db.Setup(env.ConfigEnv.Db.Url, env.ConfigEnv.Db.Name)
	log.Println(argsWithProg[1])
	chainRepo := repositories.NewChainRepository()
	chain, err := chainRepo.GetChainByName(argsWithProg[1])
	if err != nil {
		panic(err)
	}
	wsRpcs := strings.Split(chain.WsRpc, ";")
	client, _ := ethclient.DialContext(context.Background(), _chain.GetRpc(wsRpcs))

	c := crawler.NewCrawlBlock(chain.Name, client, context.Background())
	c.GetBlockNumber()
}
