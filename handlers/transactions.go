package handlers

import (
	"context"
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"net/http"
	"strings"
	"time"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"name":    "Minter Explorer Gate API",
		"version": "1.3.0",
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
	pubsubServer, ok := c.MustGet("pubsub").(*pubsub.Server)
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
		txHex := strings.ToUpper(strings.TrimSpace(tx.Transaction))
		q, _ := query.New(fmt.Sprintf("tx='%s'", txHex))
		sub, err := pubsubServer.Subscribe(context.TODO(), txHex, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": 1,
					"log":  "Subscription error",
				},
			})
			return
		}
		defer pubsubServer.Unsubscribe(context.TODO(), txHex, q)

		select {
		case <-sub.Out():
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
