package myspace

import (
	"strconv"
	"strings"
)

// helpers

func MySpaceNewDataGeneric(key string, value string) MySpaceDataPair {
	return MySpaceDataPair{Key: key, Value: value}
}

func MySpaceNewDataInt(key string, value int) MySpaceDataPair {
	return MySpaceDataPair{Key: key, Value: strconv.Itoa(value)}
}

func MySpaceNewDataBigInt(key string, value int64) MySpaceDataPair {
	return MySpaceDataPair{Key: key, Value: strconv.FormatInt(value, 10)}
}

func MySpaceNewDataBoolean(key string, value bool) MySpaceDataPair {
	if value {
		return MySpaceDataPair{Key: key, Value: "1"}
	} else {
		return MySpaceDataPair{Key: key, Value: ""}
	}
}

func MySpaceBuildPackage(datapairs []MySpaceDataPair) string {
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

func MySpaceBuildInnerBody(datapairs []MySpaceDataPair) string {

	final := ""
	for i := 0; i < len(datapairs); i++ {
		final += datapairs[i].Key + "="
		final += datapairs[i].Value + "\x1c"
	}
	return final
}

func MySpaceRetrieveKeyValue(key string, packet string) string {

	decodedPacket := packet
	splits := strings.Split(decodedPacket, "\\")

	for ix := 0; ix < len(splits); ix++ {
		if splits[ix] == key {
			return splits[ix+1]
		}
	}

	return ""
}

func MySpaceIdentifyProtocolRevision(clientver string) string {
	ver, _ := strconv.Atoi(clientver)

	if ver <= 253 {
		return "MSIMv1"
	} else if ver > 253 && ver < 366 {
		return "MSIMv2"
	} else if ver > 366 && ver < 404 {
		return "MSIMv3"
	} else if ver > 404 && ver < 594 {
		return "MSIMv4"
	} else if ver > 593 && ver < 673 {
		return "MSIMv5"
	} else if ver > 673 && ver < 697 {
		return "MSIMv6"
	} else if ver > 697 && ver < 812 {
		return "MSIMv7"
	}

	return ""
}

func MySpaceEscapeString(data string) string {
	res := strings.Replace(data, "/", "/1", -1)
	res = strings.Replace(res, "\\", "\\2", -1)
	return res
}
func MySpaceUnescapeString(data string) string {
	res := strings.Replace(data, "/1", "/", -1)
	res = strings.Replace(res, "\\2", "\\", -1)
	return res
}
