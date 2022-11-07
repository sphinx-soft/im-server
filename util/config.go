package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func readJsonConfig() map[string]interface{} {
	content, err := os.ReadFile("./config.json")

	if err != nil {
		panic(err)
	}

	// Now let's unmarshall the data into `payload`
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)

	if err != nil {
		panic(err)
	}

	return payload
}

func GetRootUrl() string {
	return fmt.Sprintf("%s", readJsonConfig()["root"])
}

func GetDatabaseLogin() string {
	return fmt.Sprintf("%s", readJsonConfig()["dblogin"])
}

func GetAESKey() string {
	return fmt.Sprintf("%s", readJsonConfig()["aeskey"])
}

func GetMailDomain() string {
	return fmt.Sprintf("%s", readJsonConfig()["maildomain"])
}

func GetServiceEnabled(service string) bool {
	str := fmt.Sprintf("%s", readJsonConfig()[service])

	if str == "on" {
		return true
	} else {
		return false
	}
}
