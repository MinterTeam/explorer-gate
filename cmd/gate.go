package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/api"
	"github.com/MinterTeam/explorer-gate/v2/core"
	"github.com/MinterTeam/minter-go-sdk/v2/api/grpc_client"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"os"
	"time"
)

var Version string   // Version
var GitCommit string // Git commit
var BuildDate string // Build date
var AppName string   // Application name

var version = flag.Bool(`v`, false, `Prints current version`)

func main() {
	flag.Parse()
	if *version {
		fmt.Printf(`%s v%s Commit %s builded %s`, AppName, Version, GitCommit, BuildDate)
		os.Exit(0)
	}

	path, err := os.Getwd()

	if fileExists(path + "/.env") {
		fmt.Printf(`loading .env file: %s`, path+".env")
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}

	//Init Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	if os.Getenv("GATE_DEBUG") != "1" && os.Getenv("GATE_DEBUG") != "true" {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.WarnLevel)
	}

	contextLogger := logger.WithFields(logrus.Fields{
		"version": Version,
		"app":     "Minter Gate",
	})

	pubsubServer := pubsub.NewServer()
	err = pubsubServer.Start()
	if err != nil {
		contextLogger.Error(err)
	}

	nodeApi, err := grpc_client.New(os.Getenv("NODE_API"))
	if err != nil {
		logrus.Fatal(err)
	}

	status, err := nodeApi.Status()
	if err != nil {
		panic(err)
	}

	latestBlock := status.LatestBlockHeight
	logger.Info(fmt.Sprintf("Starting with block %d", status.LatestBlockHeight))

	gateService := core.New(nodeApi, pubsubServer, contextLogger)

	go func() {
		for {
			block, err := nodeApi.BlockExtended(latestBlock, true)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			for _, tx := range block.Transactions {
				if tx.Code != 0 {
					err := pubsubServer.PublishWithTags(context.TODO(), "FailTx", map[string]string{
						"error": fmt.Sprintf("%X", tx.Log),
					})
					if err != nil {
						logger.Error(err)
					}
					continue
				}

				b, err := hex.DecodeString(tx.RawTx)
				if err != nil {
					logger.Error(err)
					continue
				}

				txJson, err := json.Marshal(tx)
				if err != nil {
					logger.Error(err)
					continue
				}

				err = pubsubServer.PublishWithTags(context.TODO(), "NewTx", map[string]string{
					"tx":     fmt.Sprintf("%X", b),
					"txData": string(txJson),
					"height": fmt.Sprintf("%d", block.Height),
				})
				if err != nil {
					logger.Error(err)
				}
			}
			latestBlock++
			time.Sleep(1 * time.Second)
		}
	}()

	api.Run(gateService, pubsubServer)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
