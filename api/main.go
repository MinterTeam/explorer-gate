package api

import (
	"github.com/Depado/ginprom"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/MinterTeam/explorer-gate/v2/handlers/api_v1"
	"github.com/MinterTeam/explorer-gate/v2/handlers/api_v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/libs/pubsub"
	"net/http"
	"os"
)

var Version string

// Run API
func Run(gateService *core.MinterGate, pubSubServer *pubsub.Server) {
	if os.Getenv("EXPLORER_CHECK") == "true" || os.Getenv("EXPLORER_CHECK") == "1" {
		go gateService.ExplorerStatusChecker()
	}
	router := SetupRouter(gateService, pubSubServer)
	err := router.Run(":" + os.Getenv("GATE_PORT"))
	if err != nil {
		panic(err)
	}
}

//Setup router
func SetupRouter(gateService *core.MinterGate, pubSubServer *pubsub.Server) *gin.Engine {
	router := gin.Default()
	if os.Getenv("GATE_DEBUG") != "1" && os.Getenv("GATE_DEBUG") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	if os.Getenv("GATE_DEBUG") != "1" && os.Getenv("GATE_DEBUG") != "true" {
		router.Use(gin.Recovery()) // returns 500 on any code panics
	}
	router.Use(p.Instrument())
	router.Use(cors.Default())                           // CORS
	router.Use(gin.ErrorLogger())                        // print all errors
	router.Use(apiMiddleware(gateService, pubSubServer)) // init global context

	router.GET(`/`, index)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET(`/estimate/tx-commission`, api_v1.EstimateTxCommission)
		apiV1.GET(`/estimate/coin-buy`, api_v1.EstimateCoinBuy)
		apiV1.GET(`/estimate/coin-sell`, api_v1.EstimateCoinSell)
		apiV1.GET(`/nonce/:address`, api_v1.GetNonce)
		apiV1.GET(`/min-gas`, api_v1.GetMinGas)
		apiV1.POST(`/transaction/push`, api_v1.PushTransaction)
	}

	apiV2 := router.Group("/api/v2")
	{
		apiV2.GET(`/estimate_tx_commission/:tx`, api_v2.EstimateTxCommission)
		apiV2.GET(`/estimate_coin_buy`, api_v2.EstimateCoinBuy)
		apiV2.GET(`/estimate_coin_sell`, api_v2.EstimateCoinSell)
		apiV2.GET(`/nonce/:address`, api_v2.GetNonce)
		apiV2.GET(`/min_gas_price`, api_v2.GetMinGas)
		apiV2.GET(`/coin_info/:symbol`, api_v2.CoinInfo)
		apiV2.GET(`/send_transaction/:tx`, api_v2.PushTransaction)
		apiV2.POST(`/send_transaction`, api_v2.PostTransaction)
	}
	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		err := errors.GateError{
			Error:   "",
			Code:    404,
			Message: "Resource not found.",
		}
		c.JSON(http.StatusNotFound, err)
	})
	return router
}

//Add necessary services to global context
func apiMiddleware(gate *core.MinterGate, pubSubServer *pubsub.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("gate", gate)
		c.Set("pubsub", pubSubServer)
		c.Next()
	}
}

func index(c *gin.Context) {
	gate, ok := c.MustGet("gate").(*core.MinterGate)
	if !ok {
		err := errors.GateError{
			Error:   "",
			Code:    1,
			Message: "Type cast error",
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(200, gin.H{
		"name":    "Minter Gate API",
		"version": Version,
		"active":  gate.IsActive,
	})
}
