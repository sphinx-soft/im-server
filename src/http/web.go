package http

import (
	"net/http"
	"phantom/util"
	"strconv"
	"time"
)

func RunWebServer() {
	time.Sleep(1 * time.Second) // funny sleep, prevents the log from clogging up

	if util.GetServiceEnabled("msim") {
		util.Log(util.INFO, "WebAPI Handler", "Installed IM Picture Handler for MSIM")
		http.HandleFunc("/pfp/", HandlePFP)

		util.Log(util.INFO, "WebAPI Handler", "Installed Advertisment Handler for MSIM")
		http.HandleFunc("/html.ng/", CycleMySpaceAds)
		http.HandleFunc("/adopt/", CycleMySpaceAds)
	}

	if util.GetServiceEnabled("ypager") {
		util.Log(util.INFO, "WebAPI Handler", "Installed Web Auth Handler for YMSG")
		http.HandleFunc("/config/", HandleYPager)
	}

	util.Log(util.INFO, "HTTP Listener", "Listening on 0.0.0.0:%d", 80)
	err := http.ListenAndServe(":"+strconv.Itoa(80), nil)
	if err != nil {
		util.Log(util.INFO, "WebAPI -> RunWebServer", "Error setting up http server!")
		return
	}
}
