package main

import (
	"chimera/bridge"
	"chimera/network/myspace"
	"chimera/utility"
	"os"
	"os/signal"
	"syscall"

	//	"chimera/utility/database"

	"chimera/utility/configuration"
	"chimera/utility/logging"
)

func main() {
	logging.Info("Main", "Starting Chimera...")
	logging.Info("Main", "Build Info: [%s]", utility.GetBuild())

	bridge.SignOnService("MySpace", bridge.ServiceMySpace, "2.0", configuration.GetConfiguration().Services.MySpace, myspace.LogonMySpace)
	//bridge.SignOnService("MSN", bridge.ServiceMSN, "1.0", configuration.GetConfiguration().Services.MSN, msn.LogonMSN)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for sig := range c {
		logging.Info("Exit Handler", "Captured %v! Stopping Server...", sig)
		os.Exit(0)
	}
}
