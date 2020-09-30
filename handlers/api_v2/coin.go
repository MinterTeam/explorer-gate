package api_v2

import (
	"encoding/json"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func CoinInfo(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	symbol := strings.TrimSpace(c.Param(`symbol`))

	coinInfo, err := gate.CoinInfo(symbol)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coin": symbol,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
		return
	}

	resp, err := gate.NodeClient.Marshal(coinInfo)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coin": symbol,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, json.RawMessage(resp))
	}
}

func CoinInfoById(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}

	id, err := strconv.ParseUint(strings.TrimSpace(c.Param(`id`)), 10, 64)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}

	coinInfo, err := gate.CoinInfoById(id)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coin": id,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
		return
	}

	resp, err := gate.NodeClient.Marshal(coinInfo)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coin": id,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, json.RawMessage(resp))
	}
}
