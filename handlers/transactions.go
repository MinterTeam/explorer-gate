package handlers

import (
	"github.com/daniildulin/explorer-gate/errors"
	"github.com/daniildulin/explorer-gate/helpers"
	"github.com/daniildulin/explorer-gate/services/minter_gate"
	"github.com/gin-gonic/gin"
	"github.com/olebedev/emitter"
	"net/http"
	"strings"
)

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "Minter Explorer Gate API",
		"version": "0.1",
	})
}

type PushTransactionRequest struct {
	Transaction string `form:"transaction" json:"transaction" binding:"required"`
}

func PushTransaction(c *gin.Context) {

	var err error
	gate, ok := c.MustGet("gate").(*minter_gate.MinterGate)
	helpers.CheckErrBool(ok)
	ee, ok := c.MustGet("emitter").(*emitter.Emitter)
	helpers.CheckErrBool(ok)

	var tx PushTransactionRequest
	if err = c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := gate.PushTransaction(tx.Transaction)

	if err != nil {
		switch e := err.(type) {
		case *errors.NodeError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": e.Code(),
					"log":  e.Error(),
				},
			})
		case *errors.NodeTimeOutError:
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": gin.H{
					"code": e.Code(),
					"log":  e.Error(),
				},
			})
		case *errors.InsufficientFundsError:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":  e.Code(),
					"log":   e.Error(),
					"coin":  e.Coin(),
					"value": e.Value(),
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": 0,
					"log":  e.Error(),
				},
			})
		}
	} else {
		for range ee.On(strings.ToUpper(tx.Transaction)) {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"hash": &hash,
				},
			})
			break
		}
	}
}
