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

func Test_main(t *testing.T) {
	config = env.NewViperConfig()
	router := api.SetupRouter(config)
	testPushTransaction(router, t)
}

func testPushTransaction(router *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	payload := []byte(`{"transaction":"test"}`)
	req, _ := http.NewRequest("POST", "/api/v1/transaction/push", bytes.NewBuffer(payload))
	router.ServeHTTP(w, req)
	assert.NotEqual(t, 500, w.Code)
}
