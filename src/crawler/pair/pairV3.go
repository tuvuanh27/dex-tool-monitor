package pair

import (
	"context"
	pairV3 "dex-tool/src/configs/contract/pair-v3"
	"github.com/ethereum/go-ethereum/ethclient"
)

type v3Pair struct {
	pairCrawler[pairV3.PairV3]
}

func NewV3Pair(chain string, client *ethclient.Client, ctx context.Context, pairInstance *pairV3.PairV3) IPairCrawler {
	return &v3Pair{
		pairCrawler: pairCrawler[pairV3.PairV3]{
			Chain:    chain,
			Client:   client,
			Ctx:      ctx,
			Instance: pairInstance,
		},
	}
}

func (p *v3Pair) GetFactoryAddress() string {
	instance := p.pairCrawler.Instance
	factoryAddress, err := instance.Factory(nil)
	if err != nil {
		panic(err)
	}
	return factoryAddress.Hex()
}

func (p *v3Pair) GetToken() (string, string) {
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
