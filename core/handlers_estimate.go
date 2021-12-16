package core

import (
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func (mg *MinterGate) EstimateCoinSellAllHandler(c *gin.Context) {

	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))
	coinIdToSell := strings.TrimSpace(c.Query(`coin_id_to_sell`))
	coinIdToBuy := strings.TrimSpace(c.Query(`coin_id_to_buy`))
	gasPrice := strings.TrimSpace(c.Query(`gas_price`))
	value := strings.TrimSpace(c.Query(`value_to_sell`))

	if gasPrice == "" {
		gasPrice = "1"
	}

	swapFrom := strings.TrimSpace(c.Query(`swap_from`))
	if swapFrom == "" {
		swapFrom = "optimal"
	}

	var route []uint64
	routes := c.QueryArray(`route`)
	if len(routes) == 0 {
		routes = c.QueryArray(`route[]`)
	}
	for _, r := range routes {
		cId, err := strconv.ParseUint(r, 10, 64)
		if err != nil {
			mg.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
		route = append(route, cId)
	}

	estimate, err := mg.EstimateCoinSellAll(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value, gasPrice, swapFrom, route)
	if err != nil {
		mg.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"will_get":  estimate.Value,
			"swap_from": estimate.SwapFrom,
		})
	}
}
