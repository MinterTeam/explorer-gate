package main

import (
	"flag"
	"fmt"
	"github.com/daniildulin/explorer-gate/api"
	"github.com/daniildulin/explorer-gate/env"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	api.Run(config)
}
