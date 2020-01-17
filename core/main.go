package core

import (
	"github.com/MinterTeam/minter-go-sdk/api"
	"github.com/MinterTeam/minter-go-sdk/transaction"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"os"
	"strings"
)

type MinterGate struct {
	api     *api.Api
	emitter *pubsub.Server
	Logger  *logrus.Entry
}

type CoinEstimate struct {
	Value      string
	Commission string
}

//New instance of Minter Gate
func New(e *pubsub.Server, logger *logrus.Entry) *MinterGate {
	return &MinterGate{
		emitter: e,
		api:     api.NewApi(os.Getenv("NODE_URL")),
		Logger:  logger,
	}
}

//Send transaction to blockchain
//Return transaction hash
func (mg *MinterGate) TxPush(tx string) (*string, error) {
	transactionObject, err := transaction.Decode("0x" + tx)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Warn(err)
		return nil, err
	}
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
	result, err := mg.api.EstimateCoinBuy(coinToSell, value, coinToBuy)
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
	result, err := mg.api.EstimateCoinSell(coinToSell, value, coinToBuy)
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

//Return nonce for address
func (mg *MinterGate) GetNonce(address string) (uint64, error) {
	nonce, err := mg.api.Nonce(address)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"address": address,
		}).Warn(err)
		return 0, err
	}
	return nonce - 1, nil
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
