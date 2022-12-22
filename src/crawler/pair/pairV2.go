package pair

import (
	"context"
	pairV2 "dex-tool/src/configs/contract/pair-v2"
	"github.com/ethereum/go-ethereum/ethclient"
)

type v2Pair struct {
	pairCrawler[pairV2.PairV2]
}

func NewV2Pair(chain string, client *ethclient.Client, ctx context.Context, pairInstance *pairV2.PairV2) IPairCrawler {
	return &v2Pair{
		pairCrawler: pairCrawler[pairV2.PairV2]{
			Chain:    chain,
			Client:   client,
			Ctx:      ctx,
			Instance: pairInstance,
		},
	}
}

func (p *v2Pair) GetFactoryAddress() string {
	instance := p.pairCrawler.Instance

	factoryAddress, err := instance.Factory(nil)
	if err != nil {
		panic(err)
	}
	return factoryAddress.Hex()
}

func (p *v2Pair) GetToken() (string, string) {
	instance := p.pairCrawler.Instance

	token0Address, err := instance.Token0(nil)
	if err != nil {
		panic(err)
	}

	token1Address, err := instance.Token1(nil)
	if err != nil {
		panic(err)
	}
	return token0Address.Hex(), token1Address.Hex()
}
