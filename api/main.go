package api

import (
	"context"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/daniildulin/explorer-gate/core"
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/handlers"
	"github.com/daniildulin/explorer-gate/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Run(config env.Config) {
	router := SetupRouter(config)
	router.Run(config.GetString(`gateApi.link`) + `:` + config.GetString(`gateApi.port`))
}

func SetupRouter(config env.Config) *gin.Engine {
	ee := &emitter.Emitter{}

	fmt.Println(config.GetString(`database.url`))

	db, err := gorm.Open("postgres", config.GetString(`database.url`))
	helpers.CheckErr(err)
	defer db.Close()
	db.LogMode(config.GetBool(`debug`))

	nodeRpc := tmClient.NewHTTP(`tcp://`+config.GetString(`minterApi.link`)+`:26657`, `/websocket`)
	err = nodeRpc.Start()
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q := query.MustParse(`tm.event = 'Tx'`)
	txs := make(chan interface{})

	err = nodeRpc.Subscribe(ctx, "explorer-gate", q, txs)
	if err != nil {
		log.Println(err)
	}
	go txsStore(txs, ee)

	gateService := core.New(config, ee)

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
	router.Use(gin.ErrorLogger())              // print all errors
	router.Use(gin.Recovery())                 // returns 500 on any code panics
	router.Use(apiMiddleware(gateService, ee)) // init global context

	router.GET(`/`, handlers.Index)

	v1 := router.Group("/api/v1")
	{
		v1.GET(`/estimate/tx-commission`, handlers.EstimateTxCommission)
		v1.GET(`/estimate/coin-buy`, handlers.EstimateCoinBuy)
		v1.GET(`/estimate/coin-sell`, handlers.EstimateCoinSell)

		v1.POST("/transaction/push", handlers.PushTransaction)
	}

	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found."})
	})

	return router
}

func apiMiddleware(gate *core.MinterGate, ee *emitter.Emitter) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("gate", gate)
		c.Set("emitter", ee)
		c.Next()
	}
}

func txsStore(txs <-chan interface{}, emitter *emitter.Emitter) {
	var re = regexp.MustCompile(`(?mi)^Tx\{(.*)\}`)
	for e := range txs {
		matches := re.FindStringSubmatch(e.(types.EventDataTx).Tx.String())
		<-emitter.Emit(strings.ToUpper(matches[1]), matches[1])
	}
}
