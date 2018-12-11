package minter_api

import (
	"bytes"
	"encoding/json"
	er "errors"
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/errors"
	"github.com/daniildulin/explorer-gate/helpers"
	"github.com/daniildulin/explorer-gate/models"
	"github.com/jinzhu/gorm"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

type MinterApi struct {
	config     env.Config
	db         *gorm.DB
	nodes      []models.MinterNode
	httpClient *http.Client
}

func New(config env.Config, db *gorm.DB, httpClient *http.Client) *MinterApi {
	api := &MinterApi{
		config:     config,
		db:         db,
		httpClient: httpClient,
	}
	api.GetActualNodes()
	return api
}

func (api *MinterApi) GetActiveNodesCount() int {
	return len(api.nodes)
}

func (api *MinterApi) GetActualNodes() {
	var nodes []models.MinterNode
	api.db.Where("is_excluded <> ? AND is_active = ?", true, true).Order("ping asc").Find(&nodes)
	api.nodes = nodes
}

func (api *MinterApi) PushTransaction(tx string) (string, error) {
	var err error
	response := SendTransactionResponse{}
	api.checkNodes()

	if api.GetActiveNodesCount() == 0 {
		return ``, errors.NewNodeError(`Nodes unavailable`, 0)
	}

	for _, node := range api.nodes {
		link := node.GetFullLink() + `/send_transaction?tx=0x` + tx
		err := api.getJson(link, &response)
		helpers.CheckErr(err)

		if response.Result != nil && response.Result.Code == 0 {
			return `Mt` + strings.ToLower(response.Result.Hash), nil
		} else if response.Error != nil || (response.Result != nil && response.Result.Code != 0) {
			return ``, getNodeErrorFromResponse(&response)
		}
	}

	return ``, err
}

func (api *MinterApi) GetTransaction(hash string) (bool, error) {
	var err error
	response := TransactionResponse{}
	api.checkNodes()

	if api.GetActiveNodesCount() == 0 {
		return false, errors.NewNodeError(`Nodes unavailable`, 0)
	}

	for _, node := range api.nodes {
		link := node.GetFullLink() + `/api/transaction/` + hash

		err := api.getJson(link, &response)

		if err == nil {
			return true, nil
		}
	}

	return false, err
}

func (api *MinterApi) getJson(url string, target interface{}) error {
	r, err := api.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return er.New("response code is not 200")
	}

	return json.NewDecoder(r.Body).Decode(target)
}

func (api *MinterApi) postJson(url string, data []byte, target interface{}) error {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-Minter-Network-Id", "odin")
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (api *MinterApi) checkNodes() {
	if len(api.nodes) == 0 {
		api.GetActualNodes()
	}
}

func getNodeErrorFromResponse(r *SendTransactionResponse) error {

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
