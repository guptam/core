// This is custom goose binary to support .go migration files in ./db dir

package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/onepanelio/core/db"
	v1 "github.com/onepanelio/core/pkg"
	"log"
	"os"

	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	kubeConfig := v1.NewConfig()
	client, err := v1.NewClient(kubeConfig, nil, nil)
	if err != nil {
		log.Fatalf("Failed to connect to Kubernetes cluster: %v", err)
	}
	config, err := client.GetSystemConfig()
	if err != nil {
		log.Fatalf("Failed to get system config: %v", err)
	}

	databaseDataSourceName := fmt.Sprintf("host=%v user=%v password=%v dbname=%v sslmode=disable",
		config["databaseHost"], config["databaseUsername"], config["databasePassword"], config["databaseName"])

	db := sqlx.MustConnect(config["databaseDriverName"], databaseDataSourceName)

	command := args[0]

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	if err := goose.Run(command, db.DB, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
