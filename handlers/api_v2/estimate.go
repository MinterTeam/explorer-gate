package api_v2

import (
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
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
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
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))
	value := strings.TrimSpace(c.Query(`value_to_buy`))
	estimate, err := gate.EstimateCoinBuy(coinToSell, coinToBuy, value)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"commission": estimate.Commission,
			"will_pay":   estimate.Value,
		})
	}
}

func EstimateCoinSell(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coin_to_sell`))
	coinToBuy := strings.TrimSpace(c.Query(`coin_to_buy`))
	value := strings.TrimSpace(c.Query(`value_to_sell`))
	estimate, err := gate.EstimateCoinSell(coinToSell, coinToBuy, value)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"commission": estimate.Commission,
			"will_get":   estimate.Value,
		})
	}
}

func GetNonce(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
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
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	gas, err := gate.GetMinGas()
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"gas": gas,
		})
	}
}
