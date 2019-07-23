package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/MinterTeam/explorer-gate/api"
	"github.com/MinterTeam/explorer-gate/core"
	"github.com/MinterTeam/explorer-gate/env"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"github.com/MinterTeam/minter-node-go-api"
	"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"os"
	"strconv"
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

	pubsubServer := pubsub.NewServer()
	err = pubsubServer.Start()
	if err != nil {
		contextLogger.Error(err)
	}

	gateService := core.New(config, pubsubServer, contextLogger)

	proto := `http`
	if config.GetBool(`minterApi.isSecure`) {
		proto = `https`
	}
	apiLink := proto + `://` + config.GetString(`minterApi.link`) + `:` + config.GetString(`minterApi.port`)

	nodeApi := minter_node_go_api.New(apiLink)

	latestBlockResponse, err := nodeApi.GetStatus()
	if err != nil {
		panic(err)
	}

	latestBlock, err := strconv.Atoi(latestBlockResponse.Result.LatestBlockHeight)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting with block " + strconv.Itoa(latestBlock))

	go func() {
		for {
			block, err := nodeApi.GetBlock(uint64(latestBlock))
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if block.Error != nil {
				logger.Error(block.Error.Message)
				time.Sleep(time.Second)
				continue
			}

			for _, tx := range block.Result.Transactions {
				b, _ := hex.DecodeString(tx.RawTx)
				err := pubsubServer.PublishWithTags(context.TODO(), "NewTx", map[string]string{
					"tx": fmt.Sprintf("%X", b),
				})
				if err != nil {
					logger.Error(err)
				}
			}

			latestBlock++
			logger.Info("Block " + strconv.Itoa(latestBlock))

			time.Sleep(1 * time.Second)
		}
	}()

	api.Run(config, gateService, pubsubServer)
}
