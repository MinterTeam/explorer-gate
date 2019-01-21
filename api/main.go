package api

import (
	"github.com/Depado/ginprom"
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/env"
	"github.com/MinterTeam/explorer-gate/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"net/http"
)

// Run API
func Run(config env.Config, gateService *core.MinterGate, ee *emitter.Emitter, db *gorm.DB) {
	router := SetupRouter(config, gateService, ee, db)
	router.Run(config.GetString(`gateApi.link`) + `:` + config.GetString(`gateApi.port`))
}

//Setup router
func SetupRouter(config env.Config, gateService *core.MinterGate, ee *emitter.Emitter, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	if !config.GetBool(`debug`) {
		gin.SetMode(gin.ReleaseMode)
		router.Use(gin.Logger())
	}

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())
	router.Use(gin.ErrorLogger())                  // print all errors
	router.Use(gin.Recovery())                     // returns 500 on any code panics
	router.Use(apiMiddleware(gateService, ee, db)) // init global context

	router.GET(`/`, handlers.Index)

	v1 := router.Group("/api/v1")
	{
		v1.GET(`/estimate/tx-commission`, handlers.EstimateTxCommission)
		v1.GET(`/estimate/coin-buy`, handlers.EstimateCoinBuy)
		v1.GET(`/estimate/coin-sell`, handlers.EstimateCoinSell)
		v1.GET(`/nonce/:address`, handlers.GetNonce)
		v1.GET(`/max-gas`, handlers.GetMaxGas)
		v1.POST(`/transaction/push`, handlers.PushTransaction)
	}
	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": 404, "log": "Resource not found."}})
	})
	return router
}

//Add necessary services to global context
func apiMiddleware(gate *core.MinterGate, ee *emitter.Emitter, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("gate", gate)
		c.Set("emitter", ee)
		c.Set("db", db)
		c.Next()
	}
}
