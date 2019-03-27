package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/MinterTeam/explorer-gate/api"
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/env"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/olebedev/emitter"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	tmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
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
	Version = `1.3.0`

	if config.GetBool(`debug`) {
		fmt.Println(`Service RUN on DEBUG mode`)
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Printf(`%s v%s Commit %s builded %s`, AppName, Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	//Init Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	if config.GetBool(`debug`) {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.WarnLevel)
	}

	contextLogger := logger.WithFields(logrus.Fields{
		"version": "1.3.0",
		"app":     "Minter Gate",
	})

	var err error
	ee := &emitter.Emitter{}
	gateService := core.New(config, ee, contextLogger)

	//Init RPC
	nodeRpc := tmClient.NewHTTP(`tcp://`+config.GetString(`minterApi.link`)+`:26657`, `/websocket`)
	err = nodeRpc.Start()
	if err != nil {
		contextLogger.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q := query.MustParse(`tm.event = 'Tx'`)
	txs := make(chan interface{})
	err = nodeRpc.Subscribe(ctx, "explorer-gate", q, txs)
	if err != nil {
		contextLogger.Error(err)
	}

	go handleTxs(txs, ee)

	api.Run(config, gateService, ee)
}

func handleTxs(txs <-chan interface{}, emitter *emitter.Emitter) {
	var re = regexp.MustCompile(`(?mi)^Tx\{(.*)\}`)
	for e := range txs {
		matches := re.FindStringSubmatch(e.(types.EventDataTx).Tx.String())
		<-emitter.Emit(strings.ToUpper(matches[1]), matches[1])
	}
}
