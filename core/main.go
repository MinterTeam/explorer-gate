package core

import (
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/MinterTeam/minter-go-sdk/api"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"strings"
)

type MinterGate struct {
	api     *api.Api
	Config  env.Config
	emitter *pubsub.Server
	Logger  *logrus.Entry
}

type CoinEstimate struct {
	Value      string
	Commission string
}

//New instance of Minter Gate
func New(config env.Config, e *pubsub.Server, logger *logrus.Entry) *MinterGate {

	proto := `http`
	if config.GetBool(`minterApi.isSecure`) {
		proto = `https`
	}
	apiLink := proto + `://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)
	return &MinterGate{
		emitter: e,
		api:     api.NewApi(apiLink),
		Config:  config,
		Logger:  logger,
	}
}

//Send transaction to blockchain
//Return transaction hash
func (mg *MinterGate) TxPush(tx string) (*string, error) {
	transactionObject, _ := transaction.Decode(tx)
	result, err := mg.api.SendTransaction(transactionObject)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Warn(err)
		return nil, err
	}
	hash := `Mt` + strings.ToLower(result.Hash)
	return &hash, nil
}

//Return estimate of transaction
func (mg *MinterGate) EstimateTxCommission(tx string) (*string, error) {
	transactionObject, _ := transaction.Decode(tx)
	result, err := mg.api.EstimateTxCommission(transactionObject)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Warn(err)
		return nil, err
	}
	return &result.Commission, nil
}

//Return estimate of buy coin
func (mg *MinterGate) EstimateCoinBuy(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	result, err := mg.api.EstimateCoinBuy(coinToSell, value, coinToBuy, 0)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		return nil, err
	}

	return &CoinEstimate{result.WillPay, result.Commission}, nil
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSell(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	result, err := mg.api.EstimateCoinSell(coinToSell, value, coinToBuy, 0)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		return nil, err
	}

	return &CoinEstimate{result.WillGet, result.Commission}, nil
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSellAll(coinToSell string, coinToBuy string, value string, gasPrice string) (*CoinEstimate, error) {
	//response, err := mg.api.EstimateCoinSellAll(coinToSell, coinToBuy, value, gasPrice)
	//if err != nil {
	//	mg.Logger.WithFields(logrus.Fields{
	//		"coinToSell": coinToSell,
	//		"coinToBuy":  coinToBuy,
	//		"value":      value,
	//		"gasPrice":   gasPrice,
	//	}).Warn(err)
	//	return nil, err
	//}
	//if response.Error != nil {
	//	err = errors.NewNodeError(response.Error.Message, response.Error.Code)
	//	mg.Logger.WithFields(logrus.Fields{
	//		"coinToSell": coinToSell,
	//		"coinToBuy":  coinToBuy,
	//		"value":      value,
	//		"gasPrice":   gasPrice,
	//	}).Warn(err)
	//	return nil, err
	//}
	return &CoinEstimate{"0", ""}, nil
}

//Return nonce for address
func (mg *MinterGate) GetNonce(address string) (uint64, error) {
	nonce, err := mg.api.Nonce(address)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return 0, err
	}
	return nonce, nil
}

//Return nonce for address
func (mg *MinterGate) GetMinGas() (*string, error) {
	gasPrice, err := mg.api.MinGasPrice()
	if err != nil {
		mg.Logger.Error(err)
		return nil, err
	}

	return &gasPrice, nil
}
