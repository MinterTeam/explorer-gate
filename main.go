package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/daniildulin/explorer-gate/env"
	"github.com/daniildulin/explorer-gate/handlers"
	"github.com/daniildulin/explorer-gate/helpers"
	"github.com/daniildulin/explorer-gate/services/minter_api"
	"github.com/daniildulin/explorer-gate/services/minter_gate"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var Version string   // Version
var GitCommit string // Git commit
var BuildDate string // Build date
var AppName string   // Application name
var config env.Config

var version = flag.Bool(`v`, false, `Prints current version`)

// Initialize app.
func init() {
	config = env.NewViperConfig()
	AppName = config.GetString(`name`)
	Version = `0.1`

	if config.GetBool(`debug`) {
		fmt.Println(`Service RUN on DEBUG mode`)
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf(`%s v%s Commit %s builded %s\n`, AppName, Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	ee := &emitter.Emitter{}

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

	minterApi := minter_api.New(config, db, &http.Client{Timeout: 10 * time.Second})
	gate := minter_gate.New(config, minterApi, ee)

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.ErrorLogger())       // print all errors
	router.Use(gin.Recovery())          // returns 500 on any code panics
	router.Use(apiMiddleware(gate, ee)) // init global context

	v1 := router.Group("/api/v1")

	{
		v1.POST("/transaction/push", handlers.PushTransaction)
	}

	// Default handler 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found."})
	})

	router.Run(config.GetString(`gateApi.link`) + `:` + config.GetString(`gateApi.port`))
}

func apiMiddleware(gate *minter_gate.MinterGate, ee *emitter.Emitter) gin.HandlerFunc {
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
