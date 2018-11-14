package minter_gate

import (
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/services/minter_api"
)

type MinterGate struct {
	api    *minter_api.MinterApi
	config env.Config
}

func New(config env.Config, api *minter_api.MinterApi) *MinterGate {
	return &MinterGate{
		api:    api,
		config: config,
	}
}

func (mg MinterGate) PushTransaction(transaction string) (*string, error) {
	hash, err := mg.api.PushTransaction(transaction)
	if err != nil {
		return nil, err
	}
	return &hash, nil
}
