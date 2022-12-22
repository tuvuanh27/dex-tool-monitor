package token

import (
	"context"
	erc20 "dex-tool/src/configs/contract/erc-20"
	"dex-tool/src/models/entities"
	"dex-tool/src/models/repositories"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ITokenCrawler interface {
	GetName() string
	GetSymbol() string
	GetDecimals() uint8
	GetToken(tokenRepo repositories.ITokenRepository) *entities.Token
}

type tokenCrawler struct {
	Client   *ethclient.Client
	Ctx      context.Context
	Instance *erc20.Erc20
	Address  string
}

func NewToken(client *ethclient.Client, ctx context.Context, tokenAddress string) ITokenCrawler {
	instance, err := erc20.NewErc20(common.HexToAddress(tokenAddress), client)
	if err != nil {
		panic(err)
	}
	return &tokenCrawler{
		Client:   client,
		Ctx:      ctx,
		Instance: instance,
		Address:  tokenAddress,
	}
}

func (t *tokenCrawler) GetName() string {
	instance := t.Instance
	name, err := instance.Name(nil)
	if err != nil {
		panic(err)
	}
	return name
}

func (t *tokenCrawler) GetSymbol() string {
	instance := t.Instance
	symbol, err := instance.Symbol(nil)
	if err != nil {
		panic(err)
	}
	return symbol
}

func (t *tokenCrawler) GetDecimals() uint8 {
	instance := t.Instance
	decimals, err := instance.Decimals(nil)
	if err != nil {
		panic(err)
	}
	return decimals
}

func (t *tokenCrawler) GetToken(tokenRepo repositories.ITokenRepository) *entities.Token {
	token, err := tokenRepo.GetTokenByAddress(t.Address)
	if token == nil {
		newToken := &entities.Token{
			Address:  t.Address,
			Name:     t.GetName(),
			Symbol:   t.GetSymbol(),
			Decimals: int(t.GetDecimals()),
		}
		if _, err := tokenRepo.CreateToken(newToken); err != nil {
			panic(err)
		}
		return newToken
	} else if err != nil {
		panic(err)
	}

	return token
}
