package util

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreateListener(port int) net.Listener {

	tcpServer, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(port))

	if err != nil {
		Error("TCP Listener", "Failed to start listener! %s", err.Error())
		os.Exit(1)
	}

	Log("TCP Listener", "Listening on 0.0.0.0:%d", port)

	return tcpServer
}

func WriteTrafficEx(client net.Conn, data []byte) error {

	Debug("TCP -> WriteTrafficEx", "Writing Data: %s", strings.Replace(string(data), "\r\n", "", -1))
	_, err := client.Write(data)
	return err
}

func WriteTraffic(client net.Conn, data string) error {

	Debug("TCP -> WriteTraffic", "Writing Data: %s", strings.Replace(data, "\r\n", "", -1))
	_, err := client.Write([]byte(data))
	return err
}

func ReadTraffic(client net.Conn) (data []byte, success bool) {

	client.SetReadDeadline(time.Time{})
	buf := make([]byte, 4096)
	_, err := client.Read(buf)

	if err != nil {
		Debug("TCP -> ReadTraffic", "Failed to read client traffic data: %s", err.Error())
		return buf, false
	}

	Debug("TCP -> ReadTraffic", "Reading Data: %s", strings.Replace(string(buf), "\r\n", "", -1))

	return buf, true
}

func ReadTrafficEx(client net.Conn) (data []byte, success bool) {

	client.SetReadDeadline(time.Now().Add(time.Millisecond * 300))
	buf := make([]byte, 4096)
	_, err := client.Read(buf)

	if err != nil && !strings.Contains(err.Error(), "i/o timeout") {
		Debug("TCP -> ReadTrafficEx", "Failed to read client traffic data: %s", err.Error())
		return buf, false
	} else if err != nil && strings.Contains(err.Error(), "i/o timeout") {
		buf = nil
	}

	if len(buf) > 1 {
		//strings.Replace(text, "\n", "", -1)
		Debug("TCP -> ReadTrafficEx", "Reading Data: %s", strings.Replace(string(buf), "\r\n", "", -1))
	}

	return buf, true
}
