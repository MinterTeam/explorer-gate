package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MinterTeam/explorer-gate/api"
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	db          *gorm.DB
	ee          *emitter.Emitter
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
	cfg = env.NewViperConfig()
	gateService = core.New(cfg, ee, db)
	router = api.SetupRouter(cfg, gateService, ee, db)
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
		assert.NoError(t, errors.New("Empty commission value"))
	}
}
