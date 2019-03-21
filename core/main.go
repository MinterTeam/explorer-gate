package core

import (
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/daniildulin/minter-node-api"
	"github.com/olebedev/emitter"
	"strings"
)

type MinterGate struct {
	api     *minter_node_api.MinterNodeApi
	config  env.Config
	emitter *emitter.Emitter
}

type CoinEstimate struct {
	Value      string
	Commission string
}

//New instance of Minter Gate
func New(config env.Config, e *emitter.Emitter) *MinterGate {
	proto := `http`
	if config.GetBool(`minterApi.isSecure`) {
		proto = `https`
	}
	apiLink := proto + `://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)
	return &MinterGate{
		emitter: e,
		api:     minter_node_api.New(apiLink),
		config:  config,
	}
}

//Send transaction to blockchain
//Return transaction hash
func (mg MinterGate) TxPush(transaction string) (*string, error) {
	response, err := mg.api.PushTransaction(transaction)
	if err != nil {
		return nil, err
	}
	if response.Error != nil || response.Result.Code != 0 {
		return nil, errors.GetNodeErrorFromResponse(response)
	}
	hash := `Mt` + strings.ToLower(response.Result.Hash)
	return &hash, nil
}

//Return estimate of transaction
func (mg *MinterGate) EstimateTxCommission(transaction string) (*string, error) {
	response, err := mg.api.GetEstimateTx(transaction)
	if err != nil {
		return nil, err
	}
	return &response.Result.Commission, nil
}

//Return estimate of buy coin
func (mg *MinterGate) EstimateCoinBuy(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	response, err := mg.api.GetEstimateCoinBuy(coinToSell, coinToBuy, value)
	if err != nil {
		return nil, err
	}
	return &CoinEstimate{response.Result.WillPay, response.Result.Commission}, nil
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSell(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	response, err := mg.api.GetEstimateCoinSell(coinToSell, coinToBuy, value)
	if err != nil {
		return nil, err
	}
	return &CoinEstimate{response.Result.WillGet, response.Result.Commission}, nil
}

//Return nonce for address
func (mg *MinterGate) GetNonce(address string) (*string, error) {
	response, err := mg.api.GetAddress(address)
	if err != nil {
		return nil, err
	}
	return &response.Result.TransactionCount, nil
}

//Return nonce for address
func (mg *MinterGate) GetMinGas() (*string, error) {
	response, err := mg.api.GetMinGasPrice()
	if err != nil {
		return nil, err
	}
	return &response.Result, nil
}
