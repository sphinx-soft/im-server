package msnp

import (
	"bytes"
	"fmt"
	"phantom/util"
	"strings"
)

func msnp_new_command(data string, cmd string, args string) string {
	return fmt.Sprintf("%s %s %s\r\n", cmd, getTrId(data, cmd), args)
}

func msnp_new_command_noargs(data string, cmd string) string {
	return fmt.Sprintf("%s %s\r\n", cmd, getTrId(data, cmd))
}

func findValueFromData(data_search string, packet string, offset ...int) string {

	decode := strings.Replace(packet, "\r\n", "", -1)
	splits := strings.Split(decode, " ")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == data_search {
			//return splits[ix+1+len(offset)]
			//string(bytes.Trim([]byte(splits[1]), "\x00"))
			return string(bytes.Trim([]byte(splits[ix+1+len(offset)]), "\x00"))
		}
	}

	return ""
}

func getTrId(data string, cmd string) string {
	decode := strings.Replace(data, "\r\n", "", -1)
	splits := strings.Split(decode, " ")

	return string(bytes.Trim([]byte(splits[1]), "\x00"))
}

func RemoveMsnpClient(s []*Msnp_Client, i int) []*Msnp_Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getUserData(username string) Msnp_Account {

	var acc Msnp_Account

	row, err := util.GetDatabaseHandle().Query("SELECT * from msnp WHERE email= ?", username)

	if err != nil {
		util.Error("Failed to get MSNP userdata: %s", err.Error())
	}

	row.Next()
	row.Scan(&acc.Uid, &acc.Email, &acc.Password, &acc.Screenname)
	row.Close()

	return acc
}
