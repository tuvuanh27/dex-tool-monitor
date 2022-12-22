package main

import (
	"dex-tool/src/configs/db"
	"dex-tool/src/configs/env"
	"dex-tool/src/models/repositories"
)

func init() {
	db.Setup(env.ConfigEnv.Db.Url, env.ConfigEnv.Db.Name)
}

func main() {

	//pairRepo := repositories.NewPairRepository()
	//pair, err := pairRepo.GetPairByAddress("1")
	//if err != nil {
	//	panic(err)
	//}
	//println(pair.Address)

	// get block number
	//c := crawler.NewCrawlBlock(entities.Ethereum, client, context.Background())
	//log.Println(c)
	//c.GetBlockNumber()

	//
	//pairV2, err := pair.GetPair(2, _type.Ethereum, client, context.Background(), "0xCD8B09f495A41965e22B9Dde94C71f49484986C9")
	//if err != nil {
	//	panic(err)
	//}
	//println(pairV2.GetFactoryAddress())
	//println(pairV2.GetToken())

	//redis := db.GetRedisInstance()
	//log.Println(redis.Get("DEX_TOOL_USDT_ADDRESSES_KEY"))
	//
	chainRepo := repositories.NewChainRepository()
	chainRepo.UpdateBlock("Ethereum", 2, 1000000)
	//transactionRepository := repositories.NewTransactionRepository()
	//pairRepository := repositories.NewPairRepository()

	//err := chainRepo.CreateChain(&entities.Chain{
	//	Name:    "Ethereum",
	//	HttpRpc: "https://mainnet.infura.io/v3/313ddb0fd42e4d029103e4f4cbea9c8b",
	//	WsRpc:   "wss://mainnet.infura.io/ws/v3/313ddb0fd42e4d029103e4f4cbea9c8b",
	//	TopicBlock: []entities.TopicBlock{
	//		{Block: 0, Topic: "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"},
	//		{Block: 0, Topic: "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"},
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}

	//chain, err := chainRepo.GetChainByName("Ethereum")
	//if err != nil {
	//	panic(err)
	//}
	//httpRpc := strings.Split(chain.HttpRpc, ";")
	//client, _ := ethclient.DialContext(context.Background(), _chain.GetRpc(httpRpc))
	//direc := transaction.Director{}
	//transactionV3 := &transaction.V3Transaction{}
	//direc.SetBuilder(transactionV3)
	//tx3 := direc.BuildTransaction(chain.Name, client, context.Background(), pairV3.PairV3MetaData.ABI, "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67")
	//
	//transactionV2 := &transaction.V2Transaction{}
	//direc.SetBuilder(transactionV2)
	//tx2 := direc.BuildTransaction(chain.Name, client, context.Background(), pairV2.PairV2MetaData.ABI, "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")
	//
	//log.Println(tx3.Topic)
	//log.Println(tx2.Topic)
	//_, numJobs2 := tx2.CrawlTx(16180020, 16180020)
	//_, numJobs3 := tx3.CrawlTx(16180020, 16180020)
	//log.Println(numJobs2)
	//log.Println(numJobs3)
	//
	//results := make(chan *entities.Transaction, numJobs)
	//for w := 1; w <= 3; w++ {
	//	go transactionV3.HandleTx(jobs, results)
	//}
	//
	//for a := 1; a <= numJobs; a++ {
	//	newTx := <-results
	//	_, err = transactionRepository.CreateTransaction(newTx)
	//	if err != nil {
	//		if mongo.IsDuplicateKeyError(err) {
	//			_, err := transactionRepository.UpdateTransaction(newTx)
	//			if err != nil {
	//				panic(err)
	//			}
	//		} else {
	//
	//			panic(err)
	//		}
	//	}
	//	_, err = pairRepository.UpdatePrice(newTx.Price, newTx.PriceUSD0, newTx.PriceUSD1, newTx.PairAddress)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

}
