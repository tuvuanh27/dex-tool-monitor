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
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"math"
	"strings"
)

type V3Transaction struct {
	V3Transaction Transaction
}

func (t *V3Transaction) SetChain(chain string) ITransaction {
	t.V3Transaction.Chain = chain
	return t
}

func (t *V3Transaction) SetClient(client *ethclient.Client) ITransaction {
	t.V3Transaction.Client = client
	return t
}

func (t *V3Transaction) SetCtx(ctx context.Context) ITransaction {
	t.V3Transaction.Ctx = ctx
	return t
}

func (t *V3Transaction) SetPairAbi(pairAbi string) ITransaction {
	t.V3Transaction.PairAbi = pairAbi
	return t
}

func (t *V3Transaction) SetTopic(topic string) ITransaction {
	t.V3Transaction.Topic = topic
	return t
}

func (t *V3Transaction) GetTransaction() Transaction {
	return t.V3Transaction
}

func (t *V3Transaction) CrawlTx(fromBlock, toBlock int64) (<-chan TxCrawled, int) {
	return t.V3Transaction.CrawlTx(fromBlock, toBlock)
}

func (t *V3Transaction) HandleTx(txChanel <-chan TxCrawled, results chan<- *entities.Transaction) {
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
				pairV3Crawler, _ := pairCrawler.GetPair(3, t.V3Transaction.Chain, t.V3Transaction.Client, t.V3Transaction.Ctx, i.PairAddress)
				factoryAddress := pairV3Crawler.GetFactoryAddress()
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
				token0Address, token1Address := pairV3Crawler.GetToken()
				token0Crawl := tokenCrawler.NewToken(t.V3Transaction.Client, t.V3Transaction.Ctx, token0Address)
				token1Crawl := tokenCrawler.NewToken(t.V3Transaction.Client, t.V3Transaction.Ctx, token1Address)
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
		amount0In, amount1In, amount0Out, amount1Out := getAmountTxV3(i.UnpackData[0], i.UnpackData[1])
		newTx := &entities.Transaction{
			PairAddress:      i.PairAddress,
			DexId:            pair.DexId,
			Type:             getTxType(i.UnpackData[2], i.UnpackData[3]),
			Amount0In:        amount0In,
			Amount1In:        amount1In,
			Amount0Out:       amount0Out,
			Amount1Out:       amount1Out,
			Amount0:          getAmount(i.UnpackData[0], "0"),
			Amount1:          getAmount(i.UnpackData[1], "0"),
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

		//_, err = transactionRepository.CreateTransaction(newTx)
		//if err != nil {
		//	if mongo.IsDuplicateKeyError(err) {
		//		_, err := transactionRepository.UpdateTransaction(newTx)
		//		if err != nil {
		//			panic(err)
		//		}
		//	} else {
		//
		//		panic(err)
		//	}
		//}
		//_, err = pairRepository.UpdatePrice(price, priceUsd0, priceUsd1, i.PairAddress)
		//if err != nil {
		//	panic(err)
		//}
	}
}
func getAmountTxV3(amount0, amount1 string) (amount0In, amount1In, amount0Out, amount1Out string) {
	if toInt(amount0) > 0 {
		return fmt.Sprintf("%f", math.Abs(float64(toInt(amount0)))), "0", "0", fmt.Sprintf("%f", math.Abs(float64(toInt(amount1))))
	} else {
		return "0", fmt.Sprintf("%f", math.Abs(float64(toInt(amount1)))), fmt.Sprintf("%f", math.Abs(float64(toInt(amount0)))), "0"
	}
}
