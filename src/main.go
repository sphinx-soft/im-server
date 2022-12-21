package main

import (
	"os"
	"os/signal"
	"phantom/global"
	"phantom/http"
	"phantom/msim"
	"phantom/msnp"
	"phantom/util"
	"strings"
)

func port1863Handler() {

	tcpServer := util.CreateListener(1863)

	for {
		tcpClient, err := tcpServer.Accept()

		if err != nil {
			util.Log(util.ERROR, "Port 1863 Handler", "Failed to accept Client!", err.Error())
		} else {
			util.Log(util.TRACE, "Port 1863 Handler", "Accepted Client")
		}

		var msnp_client bool
		data, success := util.ReadTrafficEx(tcpClient)

		if success {
			if strings.HasPrefix(string(data), "VER") {
				msnp_client = true
			}
		} else {
			msnp_client = false
		}

		client := global.Client{
			Connection: tcpClient,
		}

		// Handle MSNP DS Requests and redirect to 1864
		if msnp_client {
			go msnp.HandleDispatch(&client, string(data))
		} else {
			go msim.HandleClients(&client)
			go msim.HandleClientKeepalive(&client)
		}
	}
}

func main() {
	util.Log(util.INFO, "Phantom-IM", "Starting Phantom-IM-Server!")
	util.Log(util.INFO, "Phantom-IM", "Build Number: ")

	util.Log(util.INFO, "Main", "Syncing Database")
	util.InitDatabase()

	util.Log(util.INFO, "Network Service", "Server's Outbound IP: %s", util.GetOutboundIP().String())

	if util.GetServiceEnabled("msnp") || util.GetServiceEnabled("msim") {
		util.Log(util.INFO, "Service Handler", "Launched Handler for Port 1863")
		go port1863Handler()
	}

	if util.GetServiceEnabled("msnp") {

		util.Log(util.INFO, "Service Handler", "Launched Handler for MSNP Switchboard")
		go msnp.HandleSwitchboard()

		util.Log(util.INFO, "Service Handler", "Launched Handler for MSNP Notification")
		go msnp.HandleNotification()
	}

	if util.GetServiceEnabled("http") {
		util.Log(util.INFO, "Service Handler", "Launched Handler for HTTP Server")
		go http.RunWebServer(80)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		util.Log(util.INFO, "Exit Handler", "Captured %v! Stopping Server...", sig)
		os.Exit(0)
	}

}
