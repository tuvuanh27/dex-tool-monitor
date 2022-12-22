package pair

import (
	"context"
	pairV2 "dex-tool/src/configs/contract/pair-v2"
	pairV3 "dex-tool/src/configs/contract/pair-v3"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetPair(pairType int, chain string, client *ethclient.Client, ctx context.Context, pairAddress string) (IPairCrawler, error) {
	address := common.HexToAddress(pairAddress)
	switch pairType {
	case 2:
		instance, err := pairV2.NewPairV2(address, client)
		if err != nil {
			return nil, err
		}
		return NewV2Pair(chain, client, ctx, instance), nil
	case 3:
		instance, err := pairV3.NewPairV3(address, client)
		if err != nil {
			return nil, err
		}
		return NewV3Pair(chain, client, ctx, instance), nil
	default:
		return nil, errors.New("pairType incorrect")
	}
}
