package msim

import (
	"strconv"
	"strings"
)

func msim_new_data_string(key string, value string) msim_data_pair {
	return msim_data_pair{Key: key, Value: value}
}

func msim_new_data_int(key string, value int) msim_data_pair {
	return msim_data_pair{Key: key, Value: strconv.Itoa(value)}
}

func msim_new_data_int64(key string, value int64) msim_data_pair {
	return msim_data_pair{Key: key, Value: strconv.FormatInt(value, 10)}
}

func msim_new_data_dictonary(key string, value string) msim_data_pair {
	return msim_new_data_string(key, value)
}

func msim_new_data_boolean(key string, value bool) msim_data_pair {
	if value {
		return msim_data_pair{Key: key, Value: "1"}
	} else {
		return msim_data_pair{Key: key, Value: ""}
	}
}

func findValueFromKey(key string, packet []byte) string {

	decodedPacket := string(packet)
	splits := strings.Split(decodedPacket, "\\")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == key {
			return splits[ix+1]
		}
	}

	return ""
}

func buildDataPacket(datapairs []msim_data_pair) string {

	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += "\\" + datapairs[i].Key
		if datapairs[i].Value != "" {
			final += "\\" + datapairs[i].Value
		}
	}
	final += "\\final\\"
	return final
}

func buildDataBody(datapairs []msim_data_pair) string {

	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += datapairs[i].Key + "="
		final += datapairs[i].Value + "\x1c"
	}
	return final
}
