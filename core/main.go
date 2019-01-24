package core

import (
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/MinterTeam/explorer-gate/models"
	"github.com/daniildulin/minter-node-api"
	"github.com/daniildulin/minter-node-api/responses"
	"github.com/jinzhu/gorm"
	"github.com/olebedev/emitter"
	"strings"
)

type MinterGate struct {
	api     *minter_node_api.MinterNodeApi
	config  env.Config
	emitter *emitter.Emitter
	db      *gorm.DB
}

type CoinEstimate struct {
	Value      string
	Commission string
}

//New instance of Minter Gate
func New(config env.Config, e *emitter.Emitter, db *gorm.DB) *MinterGate {
	proto := `http`
	if config.GetBool(`minterApi.isSecure`) {
		proto = `https`
	}
	apiLink := proto + `://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)
	return &MinterGate{
		emitter: e,
		api:     minter_node_api.New(apiLink),
		config:  config,
		db:      db,
	}
}

//Send transaction to blockchain
//Return transaction hash
func (mg MinterGate) TxPush(transaction string) (*string, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.SendTransactionResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.PushTransaction(transaction)
			if err != nil {
				continue
			}
			if response.Error != nil || response.Result.Code != 0 {
				return nil, errors.GetNodeErrorFromResponse(response)
			}
			hash := `Mt` + strings.ToLower(response.Result.Hash)
			return &hash, nil
		}
		return nil, err
	} else {
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
}

//Return estimate of transaction
func (mg *MinterGate) EstimateTxCommission(transaction string) (*string, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.EstimateTxResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.GetEstimateTx(transaction)
			if err != nil {
				continue
			} else {
				return &response.Result.Commission, nil
			}
		}
		return nil, err
	} else {
		response, err := mg.api.GetEstimateTx(transaction)
		if err != nil {
			return nil, err
		}
		return &response.Result.Commission, nil
	}
}

//Return estimate of buy coin
func (mg *MinterGate) EstimateCoinBuy(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.EstimateCoinBuyResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.GetEstimateCoinBuy(coinToSell, coinToBuy, value)
			if err != nil {
				continue
			} else {
				return &CoinEstimate{response.Result.WillPay, response.Result.Commission}, nil
			}
		}
		return nil, err
	} else {
		response, err := mg.api.GetEstimateCoinBuy(coinToSell, coinToBuy, value)
		if err != nil {
			return nil, err
		}
		return &CoinEstimate{response.Result.WillPay, response.Result.Commission}, nil
	}
}

//Return estimate of sell coin
func (mg *MinterGate) EstimateCoinSell(coinToSell string, coinToBuy string, value string) (*CoinEstimate, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.EstimateCoinSellResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.GetEstimateCoinSell(coinToSell, coinToBuy, value)
			if err != nil {
				continue
			} else {
				return &CoinEstimate{response.Result.WillGet, response.Result.Commission}, nil
			}
		}
		return nil, err
	} else {
		response, err := mg.api.GetEstimateCoinSell(coinToSell, coinToBuy, value)
		if err != nil {
			return nil, err
		}
		return &CoinEstimate{response.Result.WillGet, response.Result.Commission}, nil
	}
}

//Return nonce for address
func (mg *MinterGate) GetNonce(address string) (*string, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.AddressResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.GetAddress(address)
			if err != nil {
				continue
			} else {
				return &response.Result.TransactionCount, nil
			}
		}
		return nil, err
	} else {
		response, err := mg.api.GetAddress(address)
		if err != nil {
			return nil, err
		}
		return &response.Result.TransactionCount, nil
	}
}

//Return nonce for address
func (mg *MinterGate) GetMinGas() (*string, error) {
	if !mg.config.GetBool(`singleNode`) && mg.db != nil {
		var err error
		var response *responses.GasResponse
		nodes := mg.GetActiveNodes()
		for _, node := range nodes {
			mg.api.SetLink(node.GetFullLink())
			response, err = mg.api.GetMinGasPrice()
			if err != nil {
				continue
			} else {
				return &response.Result, nil
			}
		}
		return nil, err
	} else {
		response, err := mg.api.GetMinGasPrice()
		if err != nil {
			return nil, err
		}
		return &response.Result, nil
	}
}

func (mg *MinterGate) GetActiveNodes() []models.MinterNode {
	var nodes []models.MinterNode
	mg.db.Where(`is_active = ?`, true).Find(&nodes)
	return nodes
}
