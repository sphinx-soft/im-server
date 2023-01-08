package tcp

import (
	"chimera/utility"
	"chimera/utility/logging"
	"fmt"
	"net"
	"os"
	"time"
)

type TcpConnection struct {
	server net.Listener
	client net.Conn
}

func CreateListener(port int) TcpConnection {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		logging.Fatal("TCP Listener", "Failed to start listener! (%s)", err.Error())
		os.Exit(0)
	}

	logging.Info("TCP Listener", "Listening on 0.0.0.0:%d", port)

	conn := TcpConnection{
		server: tcpListener,
	}

	return conn
}

func (tcp *TcpConnection) AcceptClient() error {
	lst, err := tcp.server.Accept()

	tcp.client = lst

	return err
}

func (tcp *TcpConnection) GetRemoteAddress() string {
	return tcp.client.RemoteAddr().String()
}

func (tcp *TcpConnection) WriteTraffic(data string) error {
	logging.Trace("TCP/WriteTraffic", "Writing Data: %s", utility.SanitizeString(data))
	_, err := tcp.client.Write([]byte(data))
	return err
}

func (tcp *TcpConnection) BinaryWriteTraffic(data []byte) error {
	logging.Trace("TCP/WriteTraffic", "Writing Data: %s", utility.ByteSliceToHex(data)) //untested
	_, err := tcp.client.Write([]byte(data))
	return err
}

func (tcp *TcpConnection) ReadTraffic() (data string, err error) {
	return tcp.ExReadTraffic(time.Time{})
}

func (tcp *TcpConnection) BinaryReadTraffic() (data []byte, err error) {
	return tcp.ExBinaryReadTraffic(time.Time{})
}

func (tcp *TcpConnection) ExReadTraffic(timeout time.Time) (data string, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		logging.Error("TCP/ReadTraffic", "Failed to read traffic! (%s)", err.Error())
		return "", err
	}

	ret := make([]byte, length)
	copy(ret, buf)

	retstr := utility.SanitizeString(string(ret))
	logging.Trace("TCP/ReadTraffic", "Reading Data: %s", retstr)

	return retstr, err
}

func (tcp *TcpConnection) ExBinaryReadTraffic(timeout time.Time) (data []byte, err error) {
	tcp.client.SetReadDeadline(timeout)

	buf := make([]byte, 65535)
	length, err := tcp.client.Read(buf)

	if err != nil {
		logging.Error("TCP/ReadTraffic", "Failed to read traffic! (%s)", err.Error())
		return []byte{}, err
	}

	ret := make([]byte, length)
	copy(ret, buf)

	logging.Trace("TCP/ReadTraffic", "Reading Data: %s", utility.ByteSliceToHex(ret))

	return ret, err
}

func (tcp *TcpConnection) CloseConnection() error {
	return tcp.client.Close()
}
