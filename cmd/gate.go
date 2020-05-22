package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MinterTeam/explorer-gate/v2/src/api"
	"github.com/MinterTeam/explorer-gate/v2/src/core"
	sdk "github.com/MinterTeam/minter-go-sdk/api"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/pubsub"
	"os"
	"strconv"
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

	gateService := core.New(pubsubServer, contextLogger)

	link := os.Getenv("NODE_API")
	nodeApi := sdk.NewApi(link)

	status, err := nodeApi.Status()
	if err != nil {
		panic(err)
	}

	latestBlock, err := strconv.Atoi(status.LatestBlockHeight)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting with block " + strconv.Itoa(latestBlock))

	go func() {
		for {
			block, err := nodeApi.Block(latestBlock)
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

				b, _ := hex.DecodeString(tx.RawTx)

				txJson, err := json.Marshal(tx)
				if err != nil {
					logger.Error(err)
					continue
				}

				err = pubsubServer.PublishWithTags(context.TODO(), "NewTx", map[string]string{
					"tx":     fmt.Sprintf("%X", b),
					"txData": string(txJson),
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
