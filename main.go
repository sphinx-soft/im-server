package main

import (
	"phantom/http"
	"phantom/msim"
	"phantom/util"
)

func port1863Handler() {

	tcpServer := util.CreateListener(1863)

	for {
		tcpClient, err := tcpServer.Accept()

		if err != nil {
			util.Error("Failed to accept Client! ", err.Error())
		}

		Msim := msim.Msim_Client{
			Connection: tcpClient,
			Nonce:      msim.GenerateNonce(),
		}

		go msim.HandleClients(&Msim)
		go msim.HandleClientKeepalive(&Msim)

	}
}

func main() {
	util.Log("Entry", "Starting Phantom-IM-Server!")

	util.Log("Entry", "Syncing Database")
	util.InitDatabase()

	util.Log("Handler", "Launched Handler for Port 1863")
	go port1863Handler()

	util.Log("Handler", "Launched Handler for HTTP Server")
	http.RunWebServer(80)

}
