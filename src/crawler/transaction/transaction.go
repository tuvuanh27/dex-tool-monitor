package transaction

import (
	"context"
	"dex-tool/src/models/entities"
	"dex-tool/src/models/repositories"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/big"
	"strconv"
	"strings"
)

type Transaction struct {
	Chain   string
	Client  *ethclient.Client
	Ctx     context.Context
	PairAbi string
	Topic   string
}

type TxCrawled struct {
	PairAddress      string
	BlockNumber      int64
	Timestamp        uint64
	TxHash           string
	UnpackData       []string
	TransactionIndex int
	LogIndex         int
}

func (c *Transaction) CrawlTx(fromBlock, toBlock int64) (<-chan TxCrawled, int) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(
			fromBlock),
		ToBlock: big.NewInt(
			toBlock),
		Topics: [][]common.Hash{
			{
				common.HexToHash(c.Topic),
			},
		},
	}

	logs, err := c.Client.FilterLogs(c.Ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	var blocks []uint64
	for _, v := range logs {
		blocks = append(blocks, v.BlockNumber)
	}
	timestamps := getTimestamp(unique(blocks), c.Client)
	txChanel := make(chan TxCrawled, len(logs))

	pairContract, err := abi.JSON(strings.NewReader(c.PairAbi))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("From block: ", fromBlock, " to block: ", toBlock, " total tx: ", len(logs))
	for _, vLog := range logs {
		out, err := pairContract.Unpack("Swap", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		var unPackedData []string
		for _, v := range out {
			unPackedData = append(unPackedData, v.(*big.Int).String())
		}
		var tx = TxCrawled{
			PairAddress:      vLog.Address.String(),
			BlockNumber:      int64(vLog.BlockNumber),
			Timestamp:        timestamps[vLog.BlockNumber],
			TxHash:           vLog.TxHash.String(),
			UnpackData:       unPackedData,
			TransactionIndex: int(vLog.TxIndex),
			LogIndex:         int(vLog.Index),
		}
		txChanel <- tx
	}
	close(txChanel)
	return txChanel, len(logs)
}

func getTimestamp(blockNumber []uint64, client *ethclient.Client) map[uint64]uint64 {
	res := make(map[uint64]uint64)
	for _, v := range blockNumber {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(v)))
		if err != nil {
			log.Fatal(err)
		}
		res[v] = block.Time()
	}
	return res
}

func unique(intSlice []uint64) []uint64 {
	keys := make(map[uint64]bool)
	var list []uint64
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getAmount(a, b string) string {
	x, _ := new(big.Int).SetString(a, 10)
	y, _ := new(big.Int).SetString(b, 10)
	z := big.NewInt(0).Sub(x, y)
	return z.Abs(z).String()
}

func getTxType(amount0Out, amount1Out string) entities.TransactionType {
	if toInt(amount0Out) > toInt(amount1Out) {
		return entities.Buy
	}
	return entities.Sell
}

func toInt(a string) int {
	i, _ := strconv.Atoi(a)
	return i
}

// A/USDT or USDT/A => priceUsd = price
// A/B and B in the newest pair price: get pair(has B and has price and updated_at DESC) from db(1) => has B price => A
// A/C and C not in pair has price: get pair(has A and has price and updated_at DESC) from db(1) => has A price => C
func calculatePrice(tx *entities.Transaction, pair *entities.Pair, usdtAddresses string, pairRepo repositories.IPairRepository) (price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1 float64) {
	amount0, _ := new(big.Float).SetString(tx.Amount0)
	amount1, _ := new(big.Float).SetString(tx.Amount1)
	amount0 = new(big.Float).Quo(amount0, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(pair.Token0Decimals)), nil)))
	amount1 = new(big.Float).Quo(amount1, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(pair.Token1Decimals)), nil)))
	// Case 1: A/USDT or USDT/A => priceUsd = price
	if strings.Contains(usdtAddresses, pair.Token0Address) || strings.Contains(usdtAddresses, pair.Token1Address) {

		switch true {
		case strings.Contains(usdtAddresses, pair.Token0Address):
			price, _ = new(big.Float).Quo(amount1, amount0).Float64()
			priceUsd0 = 1
			volume0, _ = amount0.Float64()
			volumeUsd0, _ = amount0.Float64()
			priceUsd1, _ = new(big.Float).Quo(amount0, amount1).Float64()
			volume1, _ = amount1.Float64()
			volumeUsd1, _ = new(big.Float).Mul(amount1, new(big.Float).SetFloat64(priceUsd1)).Float64()
			return price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1
		case strings.Contains(usdtAddresses, pair.Token1Address):
			price, _ = new(big.Float).Quo(amount1, amount0).Float64()
			priceUsd0 = price
			volume0, _ = amount0.Float64()
			volumeUsd0, _ = new(big.Float).Mul(amount0, new(big.Float).SetFloat64(priceUsd0)).Float64()
			priceUsd1 = 1
			volume1, _ = amount1.Float64()
			volumeUsd1, _ = amount1.Float64()
			return price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1
		}
	}

	// Case 2: A/B and B in the newest pair price: get pair(has B and has price and updated_at DESC) from db(1) => has B price => A
	bPrice, err := pairRepo.GetPriceToken(pair.Token1Address, pair.Address)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			log.Println("No pair has price")
		} else {
			panic(err)
		}
	}
	if bPrice != 0 {
		price, _ = new(big.Float).Quo(amount1, amount0).Float64()
		priceUsd0, _ = new(big.Float).Mul(new(big.Float).SetFloat64(price), new(big.Float).SetFloat64(bPrice)).Float64()
		volume0, _ = amount0.Float64()
		volumeUsd0, _ = new(big.Float).Mul(amount0, new(big.Float).SetFloat64(priceUsd0)).Float64()
		priceUsd1 = bPrice
		volume1, _ = amount1.Float64()
		volumeUsd1, _ = new(big.Float).Mul(amount1, new(big.Float).SetFloat64(priceUsd1)).Float64()
		return price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1
	}

	// Case 3: A/C and C not in pair has price: get pair(has A and has price and updated_at DESC) from db(1) => has A price => C
	aPrice, err := pairRepo.GetPriceToken(pair.Token0Address, pair.Address)
	if err != nil {
		if mongo.ErrNoDocuments == err {
			log.Println("No pair has price")
		} else {
			panic(err)
		}
	}
	if aPrice != 0 {
		price, _ = new(big.Float).Quo(amount1, amount0).Float64()
		priceUsd0 = aPrice
		volume0, _ = amount0.Float64()
		volumeUsd0, _ = new(big.Float).Mul(amount0, new(big.Float).SetFloat64(priceUsd0)).Float64()
		priceUsd1, _ = new(big.Float).Quo(new(big.Float).SetFloat64(priceUsd0), new(big.Float).SetFloat64(price)).Float64()
		volumeUsd1, _ = new(big.Float).Mul(amount1, new(big.Float).SetFloat64(priceUsd1)).Float64()
		return price, priceUsd0, volume0, volumeUsd0, priceUsd1, volume1, volumeUsd1
	}
	return 0, 0, 0, 0, 0, 0, 0
}
