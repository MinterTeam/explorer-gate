package api_v1

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
			Code:    "1",
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tx := strings.TrimSpace(c.Query(`transaction`))
	if tx[:2] != "0x" {
		tx = `0x` + tx
	}

	commission, err := gate.EstimateTxCommission(tx)
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"commission": &commission,
			},
		})
	}
}

func EstimateCoinBuy(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    "1",
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToBuy`))
	estimate, err := gate.EstimateCoinBuy(coinToSell, "", coinToBuy, "", value)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)
		errors.SetErrorResponse(err, c)
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
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    "1",
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToSell`))
	estimate, err := gate.EstimateCoinSell(coinToSell, "", coinToBuy, "", value)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"coinToSell": coinToSell,
			"coinToBuy":  coinToBuy,
			"value":      value,
		}).Warn(err)

		errors.SetErrorResponse(err, c)
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
		err := errors.GateError{
			Error:   "",
			Code:    "1",
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
			"data": gin.H{
				"nonce": fmt.Sprintf("%d", nonce),
			},
		})
	}
}

func GetMinGas(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    "1",
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
			"data": gin.H{
				"gas": gas,
			},
		})
	}
}
