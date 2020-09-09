package core

import (
	"github.com/MinterTeam/explorer-gate/v2/domain"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"os"
	"strconv"
	"strings"
	"time"
)

type MinterGate struct {
	nodeClient *grpc_client.Client
	emitter    *pubsub.Server
	IsActive   bool
	Logger     *logrus.Entry
}

//New instance of Minter Gate
func New(nodeApi *grpc_client.Client, e *pubsub.Server, logger *logrus.Entry) *MinterGate {
	return &MinterGate{
		emitter:    e,
		nodeClient: nodeApi,
		IsActive:   true,
		Logger:     logger,
	}
}

//Send transaction to blockchain
//Return transaction hash
func (mg *MinterGate) TxPush(tx string) (*string, error) {
	result, err := mg.nodeClient.SendTransaction(tx)
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
func (mg *MinterGate) EstimateTxCommission(tx string, optionalHeight ...int) (*string, error) {
	result, err := mg.nodeClient.EstimateTxCommission(tx, optionalHeight...)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Warn(err)
		return nil, err
	}
	return &result.Commission, nil
}

//Return estimate of buy coin
func (mg *MinterGate) EstimateCoinBuy(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value string) (*domain.CoinEstimate, error) {
	var coinToSellInfoId, coinToBuyInfoId uint64
	var err error

	if coinIdToSell != "" {
		coinToSellInfoId, err = strconv.ParseUint(coinIdToSell, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToSellInfoId, err = mg.getCoinId(coinToSell)
		if err != nil {
			return nil, err
		}
	}

	if coinIdToBuy != "" {
		coinToBuyInfoId, err = strconv.ParseUint(coinIdToBuy, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToBuyInfoId, err = mg.getCoinId(coinToBuy)
		if err != nil {
			return nil, err
		}
	}

	result, err := mg.nodeClient.EstimateCoinIDBuy(uint32(coinToSellInfoId), uint32(coinToBuyInfoId), value)
	if err != nil {
		return nil, err
	}

	return &domain.CoinEstimate{Value: result.WillPay, Commission: result.Commission}, nil
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSell(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value string) (*domain.CoinEstimate, error) {

	var coinToSellInfoId, coinToBuyInfoId uint64
	var err error

	if coinIdToSell != "" {
		coinToSellInfoId, err = strconv.ParseUint(coinIdToSell, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToSellInfoId, err = mg.getCoinId(coinToSell)
		if err != nil {
			return nil, err
		}
	}

	if coinIdToBuy != "" {
		coinToBuyInfoId, err = strconv.ParseUint(coinIdToBuy, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToBuyInfoId, err = mg.getCoinId(coinToBuy)
		if err != nil {
			return nil, err
		}
	}

	result, err := mg.nodeClient.EstimateCoinIDSell(uint32(coinToBuyInfoId), uint32(coinToSellInfoId), value)
	if err != nil {
		return nil, err
	}

	return &domain.CoinEstimate{Value: result.WillGet, Commission: result.Commission}, nil
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSellAll(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value, gasPrice string) (*domain.CoinEstimate, error) {

	var coinToSellInfoId, coinToBuyInfoId uint64
	var err error

	if coinIdToSell != "" {
		coinToSellInfoId, err = strconv.ParseUint(coinIdToSell, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToSellInfoId, err = mg.getCoinId(coinToSell)
		if err != nil {
			return nil, err
		}
	}

	if coinIdToBuy != "" {
		coinToBuyInfoId, err = strconv.ParseUint(coinIdToBuy, 10, 64)
		if err != nil {
			return nil, err
		}
	} else {
		coinToBuyInfoId, err = mg.getCoinId(coinToBuy)
		if err != nil {
			return nil, err
		}
	}

	gp, err := strconv.ParseInt(gasPrice, 10, 64)
	if err != nil {
		return nil, err
	}

	result, err := mg.nodeClient.EstimateCoinIDSellAll(uint32(coinToBuyInfoId), uint32(coinToSellInfoId), value, int(gp))
	if err != nil {
		return nil, err
	}

	return &domain.CoinEstimate{Value: result.WillGet}, nil
}

//Return nonce for address
func (mg *MinterGate) GetNonce(address string) (uint64, error) {
	nonce, err := mg.nodeClient.Nonce(address)
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
	gasPrice, err := mg.nodeClient.MinGasPrice()
	if err != nil {
		mg.Logger.Error(err)
		return nil, err
	}
	return &gasPrice.MinGasPrice, nil
}

func (mg *MinterGate) ExplorerStatusChecker() {

	sleepTime, err := strconv.ParseInt(os.Getenv("EXPLORER_CHECK_SEC"), 10, 64)
	if err != nil {
		mg.Logger.Error(err)
		return
	}
	diff, err := strconv.ParseFloat(os.Getenv("LAST_BLOCK_DIF_SEC"), 64)
	if err != nil {
		mg.Logger.Error(err)
		return
	}
	client := resty.New().SetHostURL(os.Getenv("EXPLORER_API"))

	for {
		resp, err := client.R().
			SetResult(domain.ExplorerStatusResponse{}).
			SetError(domain.ExplorerErrorResponse{}).
			Get("/api/v1/status")

		if err != nil {
			time.Sleep(time.Duration(sleepTime) * time.Second)
			continue
		}

		if resp.IsError() {
			mg.Logger.Error(resp.Error().(*domain.ExplorerErrorResponse).Error.Message)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			continue
		}

		lastBlockTime := resp.Result().(*domain.ExplorerStatusResponse).Data.LatestBlockTime
		isActive := !(time.Since(lastBlockTime).Seconds() > diff)

		if !isActive {
			mg.Logger.Error("Minter Gate is disabled")
		}
		if isActive && !mg.IsActive {
			mg.Logger.Error("Minter Gate is enabled")
		}

		mg.IsActive = isActive
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

}

func (mg *MinterGate) CoinInfo(symbol string) (*api_pb.CoinInfoResponse, error) {
	return mg.nodeClient.CoinInfo(symbol)
}

func (mg *MinterGate) getCoinId(symbol string) (uint64, error) {
	if symbol == os.Getenv("BASE_COIN") {
		return 0, nil
	}

	coinInfo, err := mg.nodeClient.CoinInfo(symbol)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(coinInfo.Id, 10, 64)
}
