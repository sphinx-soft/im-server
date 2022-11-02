package msnp

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
)

func msnp_new_command(data string, cmd string, args string) string {
	return fmt.Sprintf("%s %s %s\r\n", cmd, getTrId(data, cmd), args)
}

func msnp_new_command_noargs(data string, cmd string) string {
	return fmt.Sprintf("%s %s\r\n", cmd, getTrId(data, cmd))
}

/*
the byte trimming hack, why?:
strings.Split will run into a problem when you split with spaces
where it will stuff 4096 null bytes into your string that
you can't see without using %v and a byte cast to the string
so to fix this i added this cursed byte trimmer hack
*/

func findValueFromData(data_search string, packet string, offset int) string {

	decode := strings.Replace(packet, "\r\n", "", -1)
	splits := strings.Split(decode, " ")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == data_search {
			//return splits[ix+1+len(offset)]
			//string(bytes.Trim([]byte(splits[1]), "\x00"))
			return string(bytes.Trim([]byte(splits[ix+1+offset]), "\x00"))
		}
	}

	return ""
}

func getTrId(data string, cmd string) string {
	decode := strings.Replace(data, "\r\n", "", -1)
	splits := strings.Split(decode, " ")

	return string(bytes.Trim([]byte(splits[1]), "\x00"))
}

func generateContextKey() int {
	return rand.Intn(100000)
}

func addUserContext(ctx *msnp_context) {
	msn_context_list = append(msn_context_list, ctx)
}

func removeUserContext(s []*msnp_context, i int) []*msnp_context {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func addSwitchboardContext(ctx *msnp_switchboard_context) {
	msn_switchboard_list = append(msn_switchboard_list, ctx)
}

func removeSwitchboardContext(s []*msnp_switchboard_context, i int) []*msnp_switchboard_context {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
