package http

import (
	"encoding/base64"
	"io"
	"net/http"
	"phantom/msim"
	"phantom/util"
	"strconv"
	"strings"
)

func HandlePFP(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "=")
	image := strings.Split(split[len(split)-1], ".")
	id, err := strconv.Atoi(image[0])
	if err != nil {
		io.WriteString(w, string("err1"))
		util.Error("Error converting id to an integer")
		return
	}
	acc, err2 := msim.GetUserDataById(id)
	if err2 {
		io.WriteString(w, string("err2"))
		util.Error("error getting user object from database")
		return
	}
	res, err := base64.StdEncoding.DecodeString(acc.Avatar)
	if err != nil {
		io.WriteString(w, string("err3"))
		util.Error("Error decoding avatar")
		return
	}
	util.Debug("Provided avatar for user %s", acc.Username)
	w.Write(res)
}
