package api

import (
	"github.com/Depado/ginprom"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/explorer-gate/v2/errors"
	"github.com/MinterTeam/explorer-gate/v2/handlers/v1"
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
		apiV1.GET(`/estimate/tx-commission`, v1.EstimateTxCommission)
		apiV1.GET(`/estimate/coin-buy`, v1.EstimateCoinBuy)
		apiV1.GET(`/estimate/coin-sell`, v1.EstimateCoinSell)
		apiV1.GET(`/nonce/:address`, v1.GetNonce)
		apiV1.GET(`/min-gas`, v1.GetMinGas)
		apiV1.POST(`/transaction/push`, v1.PushTransaction)
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
