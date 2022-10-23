package util

import (
	"net"
	"os"
	"strconv"
)

func CreateListener(port int) net.Listener {

	tcpServer, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(port))

	if err != nil {
		Error("Failed to start listener! %s", err.Error())
		os.Exit(1)
	}

	Log("TCP Listener", "Listening on 0.0.0.0:%d", port)

	return tcpServer
}

func WriteTrafficEx(client net.Conn, data []byte) error {

	Log("TCP", "Writing Data: %s", string(data))
	_, err := client.Write(data)
	return err
}

func WriteTraffic(client net.Conn, data string) error {

	Log("TCP", "Writing Data: %s", data)
	_, err := client.Write([]byte(data))
	return err
}

func ReadTraffic(client net.Conn) (data []byte, success bool) {

	buf := make([]byte, 4096)
	_, err := client.Read(buf)

	if err != nil {
		Debug("Failed to read client traffic data: %s", err.Error())
		return buf, false
	}

	return buf, true
}
