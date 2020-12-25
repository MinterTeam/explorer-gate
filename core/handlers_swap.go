package core

import (
	"encoding/json"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (mg *MinterGate) SwapPoolHandler(c *gin.Context) {
	paramCoin0 := strings.TrimSpace(c.Param(`coin0`))
	paramCoin1 := strings.TrimSpace(c.Param(`coin1`))
	provider := strings.TrimSpace(c.Param(`provider`))

	coin0, err := strconv.ParseUint(paramCoin0, 10, 64)
	if err != nil {
		mg.Logger.Error(err)
		errors.SetErrorResponse(err, c)
		return
	}
	coin1, err := strconv.ParseUint(paramCoin1, 10, 64)
	if err != nil {
		mg.Logger.Error(err)
		errors.SetErrorResponse(err, c)
		return
	}

	if provider != "" {
		data, err := mg.NodeClient.SwapPoolProvider(coin0, coin1, provider)
		resp, err := mg.NodeClient.Marshal(data)
		if err != nil {
			mg.Logger.Error(err)
			errors.SetErrorResponse(err, c)
		} else {
			c.JSON(http.StatusOK, json.RawMessage(resp))
		}
	} else {
		data, err := mg.NodeClient.PairSwapPool(coin0, coin1)
		resp, err := mg.NodeClient.Marshal(data)
		if err != nil {
			mg.Logger.Warn(err)
			errors.SetErrorResponse(err, c)
		} else {
			c.JSON(http.StatusOK, json.RawMessage(resp))
		}
	}
}
