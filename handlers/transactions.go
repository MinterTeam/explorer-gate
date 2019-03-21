package handlers

import (
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/gin-gonic/gin"
	"github.com/olebedev/emitter"
	"net/http"
	"strings"
	"time"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "Minter Explorer Gate API",
		"version": "1.0",
	})
}

type PushTransactionRequest struct {
	Transaction string `form:"transaction" json:"transaction" binding:"required"`
}

func PushTransaction(c *gin.Context) {

	var err error
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
	ee, ok := c.MustGet("emitter").(*emitter.Emitter)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Type cast error",
			},
		})
		return
	}

	var tx PushTransactionRequest
	if err = c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := gate.TxPush(tx.Transaction)

	if err != nil {
		errors.SetErrorResponse(err, c)
	} else {
		select {
		case <-ee.On(strings.ToUpper(tx.Transaction), emitter.Once):
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"hash": &hash,
				},
			})
		case <-time.After(60 * time.Second):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": 1,
					"log":  `Time out waiting for transaction to be included in block`,
				},
			})
		}
	}
}
