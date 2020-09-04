package api_v2

import (
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"name":            coinInfo.Name,
			"symbol":          coinInfo.Symbol,
			"volume":          coinInfo.Volume,
			"crr":             coinInfo.Crr,
			"reserve_balance": coinInfo.ReserveBalance,
		})
	}
}
