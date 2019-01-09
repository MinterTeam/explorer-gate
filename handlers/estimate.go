package handlers

import (
	"github.com/daniildulin/explorer-gate/core"
	"github.com/daniildulin/explorer-gate/errors"
	"github.com/daniildulin/explorer-gate/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func EstimateTxCommission(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	helpers.CheckErrBool(ok)

	tx := c.Query(`transaction`)
	price, err := gate.EstimateTxCommission(tx)

	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"commission": &price,
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
