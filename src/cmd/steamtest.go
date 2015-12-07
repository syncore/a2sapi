package main

import (
	"flag"
	"fmt"
	"os"
	"steamtest/src/steam"
	"steamtest/src/steam/filters"
	"steamtest/src/util"
)

var doConfig bool

const configFlag = "config"

func init() {
	flag.BoolVar(&doConfig, configFlag, false, "Generate the configuration file")
}

func main() {
	flag.Parse()
	if doConfig {
		if err := util.CreateConfig(); err != nil {
			fmt.Printf("Unable to create configuration file: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	cfg, err := util.ReadConfig()
	//_, err := util.ReadConfig()
	if err != nil {
		fmt.Printf("Could not read configuration file '%s' in the '%s' directory.\n",
			util.ConfigFileName, util.ConfigDirectory)
		fmt.Printf("You must generate the configuration file with: %s --%s\n",
			os.Args[0], configFlag)
		os.Exit(1)
	}

	filter := filters.NewFilter(filters.GameQuakeLive, filters.SrAll, nil)
	stop := make(chan bool, 1)
	go steam.StartMasterRetrieval(stop, filter, 7, cfg.SteamConfig.TimeBetweenMasterQueries)
	<-stop

	//web.Start()

}
