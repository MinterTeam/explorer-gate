package api_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/MinterTeam/minter-go-sdk/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type PushTransactionRequest struct {
	Transaction string `form:"transaction" json:"transaction" binding:"required"`
}

func PushTransaction(c *gin.Context) {
	var err error
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	if !gate.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": 1,
				"log":  "Explorer is down",
			},
		})
		return
	}

	var timeOut int64
	timeOut, err = strconv.ParseInt(os.Getenv("NODE_API_TIMEOUT"), 10, 64)
	if err != nil {
		gate.Logger.Error(err)
		timeOut = 30 //default value
	}

	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	pubSubServer, ok := c.MustGet("pubsub").(*pubsub.Server)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}

	var tx PushTransactionRequest
	if err = c.ShouldBindJSON(&tx); err != nil {
		gate.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, errors.NewGateError(err.Error()))
		return

	}

	txn := strings.TrimSpace(tx.Transaction)
	if txn[:2] != "0x" {
		txn = `0x` + txn
	}

	_, err = gate.TxPush(txn)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"transaction": tx,
		}).Error(err)
		errors.SetErrorResponse(err, c)
	} else {
		txHex := strings.ToUpper(txn[2:])
		q, _ := query.New(fmt.Sprintf("tx='%s'", txHex))
		sub, err := pubSubServer.Subscribe(context.TODO(), txHex, q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewGateError("Subscription error"))
			return
		}
		defer pubSubServer.Unsubscribe(context.TODO(), txHex, q)

		select {
		case msg := <-sub.Out():
			gate.Logger.Log(logrus.DebugLevel, msg)
			if msg.Data() == "FailTx" {
				tags := msg.Tags()

				gate.Logger.WithFields(logrus.Fields{
					"transaction": tx,
					"code":        1,
				}).Error(tags["error"])

				c.JSON(http.StatusInternalServerError, errors.NewGateError(tags["error"]))
			} else {
				tags := msg.Tags()
				data := new(api.TransactionResult)
				err = json.Unmarshal([]byte(tags["txData"]), data)
				data.Height = tags["height"]
				c.JSON(http.StatusOK, gin.H{
					"data": gin.H{
						"hash":        data.Hash,
						"transaction": data,
					},
				})
			}
		case <-time.After(time.Duration(timeOut) * time.Second):
			gate.Logger.WithFields(logrus.Fields{
				"transaction": tx,
				"code":        "504",
			}).Error(`Time out waiting for transaction to be included in block`)

			err := errors.GateError{
				Code:    "504",
				Message: "Time out waiting for transaction to be included in block",
			}
			c.JSON(http.StatusRequestTimeout, err)

		}
	}
}
