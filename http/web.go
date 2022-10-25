package http

import (
	"net/http"
	"phantom/util"
	"strconv"
)

func RunWebServer(port int) {
	http.HandleFunc("/api", HandleAPI)
	http.HandleFunc("/pfp/", HandlePFP)
	http.HandleFunc("/html.ng/", CycleMySpaceAds)
	http.HandleFunc("/adopt/", CycleMySpaceAds)
	util.Log("HTTP Listener", "Listening on 0.0.0.0:%d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		util.Error("Error setting up http server!")
		return
	}
}
