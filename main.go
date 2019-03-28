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
	tmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"os"
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

	blocks, err := nodeRpc.Subscribe(context.TODO(), "", `tm.event = 'NewBlock'`)
	if err != nil {
		contextLogger.Error(err)
	}

	go handleBlocks(blocks, ee)

	api.Run(config, gateService, ee)
}

func handleBlocks(blocks <-chan core_types.ResultEvent, emitter *emitter.Emitter) {
	for e := range blocks {
		for _, tx := range e.Data.(types.EventDataNewBlock).Block.Txs {
			<-emitter.Emit(fmt.Sprintf("%X", []byte(tx)), true)
		}
	}
}
