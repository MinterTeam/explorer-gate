package minter_gate

import (
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/errors"
	"github.com/daniildulin/explorer-gate/services/minter_api"
	"github.com/tendermint/tendermint/types"
	"regexp"
	"strings"
	"time"
)

type MinterGate struct {
	api     *minter_api.MinterApi
	config  env.Config
	txChain <-chan interface{}
}

var maxTxTime = float64(20)

func New(config env.Config, api *minter_api.MinterApi, txs <-chan interface{}) *MinterGate {
	return &MinterGate{
		txChain: txs,
		api:     api,
		config:  config,
	}
}

func (mg MinterGate) PushTransaction(transaction string) (*string, error) {
	hash, err := mg.api.PushTransaction(transaction)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	for e := range mg.txChain {
		var re = regexp.MustCompile(`(?mi)^Tx\{(.*)\}`)
		matches := re.FindStringSubmatch(e.(types.EventDataTx).Tx.String())

		if strings.ToUpper(transaction) == matches[1] {
			return &hash, nil
		}

		if time.Since(startTime).Seconds() > maxTxTime {
			return nil, errors.NewNodeError(`Minter Gate timeout`, 504)
		}
	}

	return nil, errors.NewNodeError(`Minter Gate error`, 1)
}
