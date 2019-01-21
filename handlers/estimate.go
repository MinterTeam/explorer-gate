package handlers

import (
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/MinterTeam/explorer-gate/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func EstimateTxCommission(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	helpers.CheckErrBool(ok)
	tx := `0x` + c.Query(`transaction`)
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
	helpers.CheckErrBool(ok)
	coinToSell := c.Query(`coinToSell`)
	coinToBuy := c.Query(`coinToBuy`)
	value := c.Query(`valueToBuy`)
	estimate, err := gate.EstimateCoinBuy(coinToSell, coinToBuy, value)
	if err != nil {
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
	helpers.CheckErrBool(ok)
	coinToSell := c.Query(`coinToSell`)
	coinToBuy := c.Query(`coinToBuy`)
	value := c.Query(`valueToSell`)
	estimate, err := gate.EstimateCoinSell(coinToSell, coinToBuy, value)
	if err != nil {
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
	helpers.CheckErrBool(ok)
	address := strings.Title(c.Param(`address`))
	nonce, err := gate.GetNonce(address)
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"nonce": nonce,
			},
		})
	}
}

func GetMaxGas(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	helpers.CheckErrBool(ok)
	gas, err := gate.GetGas()
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
