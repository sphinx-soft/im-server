package msnp

import (
	"bytes"
	"fmt"
	"strings"
)

func msnp_new_command(data string, cmd string, args string) string {
	return fmt.Sprintf("%s %s %s\r\n", cmd, getTrId(data, cmd), args)
}

func msnp_new_command_noargs(data string, cmd string) string {
	return fmt.Sprintf("%s %s\r\n", cmd, getTrId(data, cmd))
}

/*
This is the worst code I've ever written!

- First offender: the byte trimming hack
strings.Split will run into a problem when you split with spaces
where it will stuff 4096 null bytes into your string that
you can't see without using %v and a byte cast to the string
so to fix this i added this cursed byte trimmer hack

- Second offender: the slice abusing trick
Golang doesn't feature preset values for params like c++
(example: void someFunc(int someInt = 0) in a header)
so to work around this i abuse a mechanism where slices
are 0 for ints when undefined thus solving my issue where
offset is not added when not defined additionally
*/

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
