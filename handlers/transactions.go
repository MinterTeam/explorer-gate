package handlers

import (
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/errors"
	"github.com/gin-gonic/gin"
	"github.com/olebedev/emitter"
	"github.com/sirupsen/logrus"
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
		gate.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := gate.TxPush(tx.Transaction)

	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Error(err)
		errors.SetErrorResponse(err, c)
	} else {
		select {
		case <-ee.On(strings.ToUpper(tx.Transaction), emitter.Once):
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"hash": &hash,
				},
			})
		case <-time.After(time.Duration(gate.Config.GetInt("minterApi.timeOut")) * time.Second):
			gate.Logger.WithFields(logrus.Fields{
				"transaction": tx,
				"code":        504,
			}).Error(`Time out waiting for transaction to be included in block`)
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": gin.H{
					"code": 1,
					"log":  `Time out waiting for transaction to be included in block`,
				},
			})
		}
	}
}
