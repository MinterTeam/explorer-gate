package api_v2

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
	Transaction string `form:"tx" json:"tx" binding:"required"`
}

func PushTransaction(c *gin.Context) {
	sendTx(strings.TrimSpace(c.Param("tx")), c)
}

func PostTransaction(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}

	var tx PushTransactionRequest
	if err := c.ShouldBindJSON(&tx); err != nil {
		gate.Logger.Error(err)
		e := errors.GateError{
			ErrorString: "",
			Code:        "1",
			Message:     err.Error(),
		}
		c.JSON(http.StatusBadRequest, e)
		return
	}

	sendTx(strings.TrimSpace(tx.Transaction), c)
}

func sendTx(tx string, c *gin.Context) {
	var err error
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Type cast error"))
		return
	}
	if !gate.IsActive {
		c.JSON(http.StatusInternalServerError, errors.NewGateError("Explorer is unavailable"))
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

	if tx[:2] != "0x" {
		tx = `0x` + tx
	}

	hash, err := gate.TxPush(tx)
	if err != nil {
		gate.Logger.WithFields(logrus.Fields{
			"tx": tx,
		}).Error(err)
		errors.SetErrorResponse(err, c)
	} else {
		txHex := strings.ToUpper(tx[2:])
		q, _ := query.New(fmt.Sprintf("tx='%s'", txHex))
		sub, err := pubSubServer.Subscribe(context.TODO(), txHex, q)
		if err != nil {
			err := errors.GateError{
				ErrorString: "",
				Code:        "1",
				Message:     "Subscription error",
			}
			c.JSON(http.StatusInternalServerError, err)
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

				c.JSON(http.StatusBadRequest, errors.NewGateError(tags["error"]))
			} else {
				tags := msg.Tags()
				data := new(api.TransactionResult)
				err = json.Unmarshal([]byte(tags["txData"]), data)
				data.Height = tags["height"]
				c.JSON(http.StatusOK, gin.H{
					"hash": &hash,
					"data": data,
					"code": data.Code,
					"log":  data.Log,
				})
			}
		case <-time.After(time.Duration(timeOut) * time.Second):
			gate.Logger.WithFields(logrus.Fields{
				"transaction": tx,
				"code":        "504",
			}).Error(`Time out waiting for transaction to be included in block`)

			err := errors.GateError{
				ErrorString: "",
				Code:        "504",
				Message:     "Time out waiting for transaction to be included in block",
			}
			c.JSON(http.StatusRequestTimeout, err)
		}
	}
}
