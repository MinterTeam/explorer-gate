package api_v2

import (
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func EstimateTxCommission(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}

	tx := strings.TrimSpace(c.Param(`tx`))
	if tx[:2] != "0x" {
		tx = `0x` + tx
	}

	commission, err := gate.EstimateTxCommission(tx)
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"commission": &commission,
		})
	}
}

func EstimateCoinBuy(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}

	swapFrom := strings.TrimSpace(c.Query(`swap_from`))
	if swapFrom == "" {
		swapFrom = "optimal"
	}

	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinIdToSell := strings.TrimSpace(c.Query(`coin_id_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))
	coinIdToBuy := strings.TrimSpace(c.Query(`coin_id_to_buy`))
	value := strings.TrimSpace(c.Query(`value_to_buy`))
	var route []uint64
	if strings.TrimSpace(c.Query(`route`)) != "" {
		err := json.Unmarshal([]byte(strings.TrimSpace(c.Query(`route`))), &route)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)

			errors.SetErrorResponse(err, c)
			return
		}
	}
	estimate, err := gate.EstimateCoinBuy(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value, swapFrom, route)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, estimate)
	}
}

func EstimateCoinSell(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}

	swapFrom := strings.TrimSpace(c.Query(`swap_from`))
	if swapFrom == "" {
		swapFrom = "optimal"
	}

	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))

	coinIdToSell := strings.TrimSpace(c.Query(`coin_id_to_sell`))
	coinIdToBuy := strings.TrimSpace(c.Query(`coin_id_to_buy`))

	value := strings.TrimSpace(c.Query(`value_to_sell`))

	var route []uint64
	if strings.TrimSpace(c.Query(`route`)) != "" {
		err := json.Unmarshal([]byte(strings.TrimSpace(c.Query(`route`))), &route)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)

			errors.SetErrorResponse(err, c)
			return
		}
	}

	estimate, err := gate.EstimateCoinSell(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value, swapFrom, route)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, estimate)
	}
}

func EstimateCoinSellAll(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))
	coinIdToSell := strings.TrimSpace(c.Query(`coin_id_to_sell`))
	coinIdToBuy := strings.TrimSpace(c.Query(`coin_id_to_buy`))
	gasPrice := strings.TrimSpace(c.Query(`gas_price`))
	value := strings.TrimSpace(c.Query(`value_to_sell`))

	swapFrom := strings.TrimSpace(c.Query(`swap_from`))
	if swapFrom == "" {
		swapFrom = "optimal"
	}

	var route []uint64
	if strings.TrimSpace(c.Query(`route`)) != "" {
		err := json.Unmarshal([]byte(strings.TrimSpace(c.Query(`route`))), &route)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)

			errors.SetErrorResponse(err, c)
			return
		}
	}

	estimate, err := gate.EstimateCoinSellAll(coinToSell, coinIdToSell, coinToBuy, coinIdToBuy, value, gasPrice, swapFrom, route)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, estimate)
	}
}

func GetNonce(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}
	address := strings.Title(strings.TrimSpace(c.Param(`address`)))
	nonce, err := gate.GetNonce(address)
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"nonce": fmt.Sprintf("%d", nonce),
		})
	}
}

func GetMinGas(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusRequestTimeout, gin.H{
			"error": errors.NewGateError("Type cast error"),
		})
		return
	}
	gas, err := gate.GetMinGas()
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"min_gas_price": gas,
		})
	}
}
