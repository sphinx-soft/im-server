package main

import (
	"phantom/utility"

	//	"phantom/utility/database"

	"phantom/utility/logging"
)

func main() {
	logging.Info("Main", "Starting Phantom...")
	logging.Info("Main", "Build Info: [%s]", utility.GetBuild())

}
