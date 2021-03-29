package api_v1

import (
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func EstimateTxCommission(c *gin.Context) {
	var err error
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}

	tx := strings.TrimSpace(c.Query(`transaction`))
	if tx[:2] != "0x" {
		tx = `0x` + tx
	}

	commission, err := gate.EstimateTxCommission(tx)
	if err != nil {
		errors.SetErrorResponseV1(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"commission": &commission,
			},
		})
	}
}

func EstimateCoinBuy(c *gin.Context) {
	var err error
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToBuy`))

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
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
		route = append(route, cId)
	}

	commissionIdCoin := uint64(0)
	commissionCoin := strings.TrimSpace(c.Query(`coin_commission`))
	if strings.TrimSpace(c.Query(`coin_id_commission`)) != "" {
		commissionIdCoin, err = strconv.ParseUint(c.Query(`coin_id_commission`), 10, 64)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
	} else if commissionCoin != "" {
		coinInfo, err := gate.CoinInfo(commissionCoin)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
		commissionIdCoin = coinInfo.Id
	}

	estimate, err := gate.EstimateCoinBuy(coinToSell, "", coinToBuy, "", value, swapFrom, commissionIdCoin, route)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		errors.SetErrorResponseV1(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"commission": estimate.Commission,
				"will_pay":   estimate.Value,
			},
		})
	}
}

func EstimateCoinSell(c *gin.Context) {
	var err error
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToSell`))

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
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
		route = append(route, cId)
	}

	commissionIdCoin := uint64(0)
	commissionCoin := strings.TrimSpace(c.Query(`coin_commission`))
	if strings.TrimSpace(c.Query(`coin_id_commission`)) != "" {
		commissionIdCoin, err = strconv.ParseUint(c.Query(`coin_id_commission`), 10, 64)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
	} else if commissionCoin != "" {
		coinInfo, err := gate.CoinInfo(commissionCoin)
		if err != nil {
			gate.Logger.WithFields(logrus.Fields{
				"coinToSell": coinToSell,
				"coinToBuy":  coinToBuy,
				"value":      value,
			}).Warn(err)
			errors.SetErrorResponse(err, c)
			return
		}
		commissionIdCoin = coinInfo.Id
	}

	estimate, err := gate.EstimateCoinSell(coinToSell, "", coinToBuy, "", value, swapFrom, commissionIdCoin, route)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponseV1(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"commission": estimate.Commission,
				"will_get":   estimate.Value,
			},
		})
	}
}

func GetNonce(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	address := strings.Title(strings.TrimSpace(c.Param(`address`)))
	nonce, err := gate.GetNonce(address)
	if err != nil {
		errors.SetErrorResponseV1(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"nonce": fmt.Sprintf("%d", nonce),
			},
		})
	}
}

func GetMinGas(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	gas, err := gate.GetMinGas()
	if err != nil {
		errors.SetErrorResponseV1(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"gas": gas,
			},
		})
	}
}
