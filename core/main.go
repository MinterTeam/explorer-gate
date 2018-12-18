package core

import (
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/errors"
	"github.com/daniildulin/minter-node-api"
	"github.com/daniildulin/minter-node-api/responses"
	"github.com/olebedev/emitter"
	"math/big"
	"regexp"
	"strings"
)

type MinterGate struct {
	api     *minter_node_api.MinterNodeApi
	config  env.Config
	emitter *emitter.Emitter
}

func New(config env.Config, e *emitter.Emitter) *MinterGate {

	apiLink := `http`
	if config.GetBool(`minterApi.isSecure`) {
		apiLink += `s://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)
	} else {
		apiLink += `://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)
	}

	return &MinterGate{
		emitter: e,
		api:     minter_node_api.New(apiLink),
		config:  config,
	}
}

func (mg MinterGate) TxPush(transaction string) (*string, error) {
	response, err := mg.api.PushTransaction(transaction)
	if err != nil {
		return nil, err
	}

	if response.Error != nil || response.Result.Code != 0 {
		return nil, getNodeErrorFromResponse(response)
	}

	hash := `Mt` + strings.ToLower(response.Result.Hash)
	return &hash, nil
}

func getNodeErrorFromResponse(r *responses.SendTransactionResponse) error {

	bip := big.NewFloat(0.000000000000000001)

	if r.Result != nil {
		switch r.Result.Code {
		case 107:
			var re = regexp.MustCompile(`(?mi)^.*Wanted *(\d+) (\w+)`)
			matches := re.FindStringSubmatch(r.Result.Log)
			value, _, err := big.ParseFloat(matches[1], 10, 0, big.ToZero)
			if err != nil {
				return err
			}
			value = value.Mul(value, bip)
			return errors.NewInsufficientFundsError(strings.Replace(r.Result.Log, matches[1], value.String(), -1), int32(r.Result.Code), value.String(), matches[2])
		default:
			return errors.NewNodeError(r.Result.Log, int32(r.Result.Code))
		}
	}

	if r.Error != nil {
		return errors.NewNodeError(r.Error.Data, r.Error.Code)
	}

	return errors.NewNodeError(`Unknown error`, -1)
}
