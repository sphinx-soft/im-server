package msnp

import (
	"fmt"
	"strings"
)

func msnp_new_command(data string, cmd string, args string) msnp_command {
	return msnp_command{command: cmd, transactionID: findValue(data, cmd), arguments: args}
}

func buildCommand(command_data msnp_command) string {
	return fmt.Sprintf("%s %s %s\r\n", command_data.command, command_data.transactionID, command_data.arguments)
}

func findValue(data string, cmd string) string {
	splits := strings.Split(data, " ")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == cmd {
			return splits[ix+1]
		}
	}

	return ""
}

func RemoveMsnpClient(s []*Msnp_Client, i int) []*Msnp_Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
