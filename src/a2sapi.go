package main

import (
	"a2sapi/src/config"
	"a2sapi/src/constants"
	"a2sapi/src/db"
	"a2sapi/src/steam"
	"a2sapi/src/steam/filters"
	"a2sapi/src/util"
	"a2sapi/src/web"
	"flag"
	"fmt"
	"os"
)

var (
	doConfig       bool
	useDebugConfig bool
	runSilent      bool
)

const (
	configFlag = "config"
	debugFlag  = "debug"
	silentFlag = "silent"
)

func init() {
	flag.BoolVar(&doConfig, configFlag, false, "Generate the configuration file")
	flag.BoolVar(&useDebugConfig, debugFlag, false, "Use debug mode configuration file")
	flag.BoolVar(&runSilent, silentFlag, false,
		"Launch without displaying startup information")
}

func main() {
	flag.Parse()

	if doConfig {
		if !util.FileExists(constants.GameFileFullPath) {
			filters.DumpDefaultGames()
		}
		config.CreateConfig()
		os.Exit(0)
	}

	if useDebugConfig {
		config.CreateDebugConfig()
		constants.IsDebug = true
		launch(true)
	} else {
		launch(false)
	}
}

func launch(isDebug bool) {
	if !util.FileExists(constants.GameFileFullPath) {
		filters.DumpDefaultGames()
	}
	if !isDebug {
		if !util.FileExists(constants.ConfigFilePath) {
			fmt.Printf("Could not read configuration file '%s' in the '%s' directory.\n",
				constants.ConfigFilename, constants.ConfigDirectory)
			fmt.Printf("You must generate the configuration file with: %s --%s\n",
				os.Args[0], configFlag)
			os.Exit(1)
		}
	}
	// Initialize the application-wide configuration
	config.InitConfig()
	// Initialize the application-wide database connections (panic on failure)
	db.InitDBs()

	if !runSilent {
		printStartInfo()
	}

	if config.Config.SteamConfig.AutoQueryMaster {
		autoQueryGame := filters.GetGameByName(
			config.Config.SteamConfig.AutoQueryGame)
		if autoQueryGame == filters.GameUnspecified {
			fmt.Println("Invalid game specified for automatic timed query!")
			fmt.Printf(
				"You may need to delete: '%s' and/or recreate the config with: %s --%s",
				constants.GameFileFullPath, os.Args[0], configFlag)
			os.Exit(1)
		}
		// HTTP server + API + Steam auto-querier
		go web.Start(runSilent)
		filter := filters.NewFilter(autoQueryGame, filters.SrAll, nil)
		stop := make(chan bool, 1)
		go steam.StartMasterRetrieval(stop, filter, 7,
			config.Config.SteamConfig.TimeBetweenMasterQueries)
		<-stop
	} else {
		// HTTP server + API standalone
		web.Start(runSilent)
	}
}

func printStartInfo() {
	fmt.Printf("%s\n", constants.AppInfo)
	if useDebugConfig {
		fmt.Println("NOTE: We're currently using debug the configuration!")
	}
	if config.Config.SteamConfig.AutoQueryMaster {
		fmt.Println("Automatic timed master server queries: enabled")
		fmt.Printf("Automatic timed master server queries every %d seconds\n",
			config.Config.SteamConfig.TimeBetweenMasterQueries)
		fmt.Printf("Automatic timed master server query game: %s\n",
			config.Config.SteamConfig.AutoQueryGame)
		fmt.Printf("Automatic timed master server query max hosts to receive: %d\n",
			config.Config.SteamConfig.MaximumHostsToReceive)
	} else {
		fmt.Println("Automatic timed master server queries: disabled")
	}
}
