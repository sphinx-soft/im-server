package main

import (
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
			if strings.Contains(string(data), "VER") {
				msnp_client = true
			}
		} else {
			msnp_client = false
		}

		// Handle MSNP DS Requests and redirect to 1864
		if msnp_client {
			Msnp := msnp.Msnp_Client{
				Connection: tcpClient,
			}

			go msnp.HandleDispatch(&Msnp, string(data))

		} else {
			Msim := msim.Msim_Client{
				Connection: tcpClient,
				Nonce:      msim.GenerateNonce(),
			}

			go msim.HandleClients(&Msim)
			go msim.HandleClientKeepalive(&Msim)

		}
	}
}

func main() {
	util.Log("Entry", "Starting Phantom-IM-Server!")

	util.Log("Entry", "Syncing Database")
	util.InitDatabase()

	util.Log("Handler", "Launched Handler for Port 1863")
	go port1863Handler()

	util.Log("Handler", "Launched Handler for HTTP Server")
	go http.RunWebServer(80)

	util.Log("Handler", "Launched Handler for MSNP Notification")
	msnp.HandleNotification()
}
