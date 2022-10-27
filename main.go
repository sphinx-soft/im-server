package main

import (
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
			util.Error("Failed to accept Client! ", err.Error())
		} else {
			util.Debug("Port 1836 Handler", "Accepted Client")
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
	util.Log("Entry", "Starting Phantom-IM-Server!")

	util.Log("Entry", "Syncing Database")
	util.InitDatabase()

	util.Log("Handler", "Launched Handler for Port 1863")
	go port1863Handler()

	util.Log("Handler", "Launched Handler for MSNP Notification")
	go msnp.HandleNotification()

	util.Log("Handler", "Launched Handler for HTTP Server")
	http.RunWebServer(80)
}
