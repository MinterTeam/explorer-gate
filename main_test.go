package main

import (
	"bytes"
	"github.com/daniildulin/explorer-gate/api"
	"github.com/daniildulin/explorer-gate/env"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testTx = `0xf8820d018a4d4e540000000000000001a9e88a4d4e5400000000000000941b685a7c1e78726c48f619c497a07ed75fe00483872386f26fc10000808001b845f8431ca05ddcd3ffd2d5b21ffe4686cadbb462bad9facdd7ee0c2db31a7b6da6f06468b3a044df8fc8b4c4190ef352e0f70112527b6b25c4a22a67c9e9365ac7e511ac12f3`

func Test_main(t *testing.T) {
	config = env.NewViperConfig()
	router := api.SetupRouter(config)
	testPushTransaction(router, t)
	testEstimateTx(router, t)
}

func testPushTransaction(router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	payload := []byte(`{"transaction":"test"}`)
	req, _ := http.NewRequest("POST", "/api/v1/transaction/push", bytes.NewBuffer(payload))
	router.ServeHTTP(w, req)
	assert.NotEqual(t, 500, w.Code)
}

func testEstimateTx(router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/estimate/tx-commission?transaction="+testTx, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
