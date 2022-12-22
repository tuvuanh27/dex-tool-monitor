package transaction

import (
	"context"
	"dex-tool/src/configs/db"
	"dex-tool/src/configs/env"
	pairCrawler "dex-tool/src/crawler/pair"
	tokenCrawler "dex-tool/src/crawler/token"
	"dex-tool/src/models/entities"
	"dex-tool/src/models/repositories"
	"encoding/json"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"strings"
)

type V2Transaction struct {
	V2Transaction Transaction
}

func (t *V2Transaction) SetChain(chain string) ITransaction {
	t.V2Transaction.Chain = chain
	return t
}

func (t *V2Transaction) SetClient(client *ethclient.Client) ITransaction {
	t.V2Transaction.Client = client
	return t
}

func (t *V2Transaction) SetCtx(ctx context.Context) ITransaction {
	t.V2Transaction.Ctx = ctx
	return t
}

func (t *V2Transaction) SetPairAbi(pairAbi string) ITransaction {
	t.V2Transaction.PairAbi = pairAbi
	return t
}

func (t *V2Transaction) SetTopic(topic string) ITransaction {
	t.V2Transaction.Topic = topic
	return t
}

func (t *V2Transaction) GetTransaction() Transaction {
	return t.V2Transaction
}

func (t *V2Transaction) CrawlTx(fromBlock, toBlock int64) (<-chan TxCrawled, int) {
	return t.V2Transaction.CrawlTx(fromBlock, toBlock)
}

func (t *V2Transaction) HandleTx(txChanel <-chan TxCrawled, results chan<- *entities.Transaction) {
	pairRepository := repositories.NewPairRepository()
	dexRepository := repositories.NewDexRepository()
	tokenRepository := repositories.NewTokenRepository()
	redisInstance := db.GetRedisInstance()
	usdtAddress := strings.Join(env.ConfigEnv.Stable, ",")

	for i := range txChanel {
		var pair *entities.Pair
		p, err := redisInstance.Get(i.PairAddress)
		if err == redis.Nil {
			pair, _ = pairRepository.GetPairByAddress(i.PairAddress)
			if pair == nil {
				var dex *entities.Dex
				pairV2Crawler, _ := pairCrawler.GetPair(2, t.V2Transaction.Chain, t.V2Transaction.Client, t.V2Transaction.Ctx, i.PairAddress)
				factoryAddress := pairV2Crawler.GetFactoryAddress()
				sDex, err := redisInstance.Get(factoryAddress)
				if err == redis.Nil {
					dex, _ = dexRepository.GetDexByFactory(factoryAddress)
					if dex == nil {
						continue
					}
				} else {
					err = json.Unmarshal([]byte(sDex), &dex)
					if err != nil {
						panic(err)
					}
				}
				token0Address, token1Address := pairV2Crawler.GetToken()
				token0Crawl := tokenCrawler.NewToken(t.V2Transaction.Client, t.V2Transaction.Ctx, token0Address)
				token1Crawl := tokenCrawler.NewToken(t.V2Transaction.Client, t.V2Transaction.Ctx, token1Address)
				token0 := token0Crawl.GetToken(tokenRepository)
				token1 := token1Crawl.GetToken(tokenRepository)
				pair = &entities.Pair{
					Address:        i.PairAddress,
					DexId:          dex.ID,
					Token0Address:  token0.Address,
					Token1Address:  token1.Address,
					Token0Decimals: token0.Decimals,
					Token1Decimals: token1.Decimals,
					Token0Symbol:   token0.Symbol,
					Token1Symbol:   token1.Symbol,
				}
				if _, err := pairRepository.CreatePair(pair); err != nil {
					panic(err)
				}

			}
			pairCache, _ := json.Marshal(pair)
			if err := redisInstance.Set(pair.Address, string(pairCache), 0); err != nil {
				panic(err)
			}

		} else {
			err = json.Unmarshal([]byte(p), &pair)
			if err != nil {
				panic(err)
			}
		}

		newTx := &entities.Transaction{
			PairAddress:      i.PairAddress,
			DexId:            pair.DexId,
			Type:             getTxType(i.UnpackData[2], i.UnpackData[3]),
			Amount0In:        i.UnpackData[0],
			Amount1In:        i.UnpackData[1],
			Amount0Out:       i.UnpackData[2],
			Amount1Out:       i.UnpackData[3],
			Amount0:          getAmount(i.UnpackData[0], i.UnpackData[2]),
			Amount1:          getAmount(i.UnpackData[1], i.UnpackData[3]),
			BlockNumber:      i.BlockNumber,
			Timestamp:        int64(i.Timestamp),
			TransactionHash:  i.TxHash,
			TransactionIndex: i.TransactionIndex,
			LogIndex:         i.LogIndex,
		}
		price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1 := calculatePrice(newTx, pair, usdtAddress, pairRepository)
		newTx.Price = price
		newTx.PriceUSD0 = priceUsd0
		newTx.Volume0 = volume0
		newTx.VolumeUSD0 = volumeUsd0
		newTx.PriceUSD1 = priceUsd1
		newTx.Volume1 = volume1
		newTx.VolumeUSD1 = volumeUsd1

		results <- newTx
	}
}
