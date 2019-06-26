package handlers

import (
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func EstimateTxCommission(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	tx := `0x` + strings.TrimSpace(c.Query(`transaction`))
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToBuy`))
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
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToSell`))
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

func EstimateCoinSellAll(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	coinToSell := strings.TrimSpace(c.Query(`coinToSell`))
	coinToBuy := strings.TrimSpace(c.Query(`coinToBuy`))
	value := strings.TrimSpace(c.Query(`valueToSell`))
	gasPrice := strings.TrimSpace(c.Query(`gasPrice`))
	estimate, err := gate.EstimateCoinSellAll(coinToSell, coinToBuy, value, gasPrice)
	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"will_get": estimate.Value,
			},
		})
	}
}

func GetNonce(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}
	address := strings.Title(strings.TrimSpace(c.Param(`address`)))
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

func GetMinGas(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
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
