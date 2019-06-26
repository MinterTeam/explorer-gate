package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MinterTeam/explorer-gate/api"
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/pubsub"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	cfg         env.Config
	router      *gin.Engine
	gateService *core.MinterGate
	testTx      = `f8820d018a4d4e540000000000000001a9e88a4d4e5400000000000000941b685a7c1e78726c48f619c497a07ed75fe00483872386f26fc10000808001b845f8431ca05ddcd3ffd2d5b21ffe4686cadbb462bad9facdd7ee0c2db31a7b6da6f06468b3a044df8fc8b4c4190ef352e0f70112527b6b25c4a22a67c9e9365ac7e511ac12f3`
)

type RespData struct {
	Commission *string `json:"commission"`
}

type RespError struct {
	Code  int     `json:"code"`
	Log   string  `json:"log"`
	Value *int    `json:"value"`
	Coin  *string `json:"coin"`
}

type Resp struct {
	Data  *RespData  `json:"data"`
	Error *RespError `json:"error"`
}

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	contextLogger := logger.WithFields(logrus.Fields{
		"version": "1.3.0",
		"app":     "Minter Gate Test",
	})

	pubsubServer := pubsub.NewServer()
	err := pubsubServer.Start()
	if err != nil {
		panic(err)
	}

	cfg = env.NewViperConfig()
	gateService = core.New(cfg, pubsubServer, contextLogger)
	router = api.SetupRouter(cfg, gateService, pubsubServer)
}

func TestPushWrongTransaction(t *testing.T) {
	var target Resp
	w := httptest.NewRecorder()
	payload := []byte(`{"transaction":"` + testTx + `"}`)
	req, err := http.NewRequest("POST", "/api/v1/transaction/push", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	err = json.NewDecoder(w.Body).Decode(&target)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.IsType(t, target.Error.Code, int(1))
	assert.IsType(t, target.Error.Log, "")
}

func TestEstimateTx(t *testing.T) {
	var target Resp
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/estimate/tx-commission?transaction="+testTx, nil)
	router.ServeHTTP(w, req)

	err := json.NewDecoder(w.Body).Decode(&target)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, err)
	assert.NotNil(t, target.Data)
	assert.NotNil(t, target.Data.Commission)
	if target.Data.Commission != nil && *target.Data.Commission == "" {
		assert.NoError(t, errors.New("empty commission value"))
	}
}

func TestErrors(t *testing.T) {
	bip := big.NewFloat(0.000000000000000001)
	zero := big.NewFloat(0)

	//min := big.NewFloat(-10000000000)
	//

	//errorString := fmt.Sprintf("You wanted to get minimum %s, but currently you will get %s", "0", "-386524121155570000000000")
	//errorString := fmt.Sprintf("Coin reserve balance is not sufficient for transaction. Has: %s, required %s")
	//errorString := fmt.Sprintf("Insufficient funds for sender account: %s. Wanted %s %s")

	//NO BIP
	//101, 102 105 109 110 113 115  603 604
	//errorString := fmt.Sprintf("Response code is not correct. Expected %d, got %d")
	//errorString := fmt.Sprintf("Response code is not correct. Expected %d, got %d")
	//
	//errorString := fmt.Sprintf()
	//
	//errorString := fmt.Sprintf()

	//604
	//errorString := fmt.Sprintf("Not enough multisig votes. Needed %d, has %d")

	//BIP
	//107
	//errorString := fmt.Sprintf("Gas coin reserve balance is not sufficient for transaction. Has: %s %s, required %s %s")
	//errorString := fmt.Sprintf("Insufficient funds for sender account: %s. Wanted %s %s")
	//114
	//errorString := fmt.Sprintf("Gas price of tx is too low to be included in mempool. Expected %d")
	//302
	//errorString := fmt.Sprintf("You wanted to sell maximum %s, but currently you need to spend %s to complete tx")
	//303
	//errorString := fmt.Sprintf("You wanted to get minimum %s, but currently you will get %s",)
	//
	//errorString := fmt.Sprintf()
	//errorString := fmt.Sprintf()
	//errorString := fmt.Sprintf()

	errorString := fmt.Sprintf("You wanted to get minimum %s, but currently you will get %s", "0", "-38652412115557")
	var re = regexp.MustCompile(`(?mi)(Has:|required|Wanted|Expected|maximum|spend|minimum|get) -*(\d+)`)
	matches := re.FindAllStringSubmatch(errorString, -1)

	if matches != nil {
		for _, match := range matches {
			var valueString string
			value, _, err := big.ParseFloat(match[2], 10, 0, big.ToZero)
			if err != nil {
				return
			}
			value = value.Mul(value, bip)

			if value.Cmp(zero) == 0 {
				valueString = "0"
			} else {
				valueString = value.Text('f', 10)
			}
			replacer := strings.NewReplacer(match[2], valueString)
			errorString = replacer.Replace(errorString)
		}
		fmt.Println(errorString)
	}
}
