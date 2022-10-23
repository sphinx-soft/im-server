package http

import (
	"io"
	"net/http"
	"phantom/util"
)

func CycleMySpaceAds(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	util.Log("WEB", "Sending MySpace Advertisment Data!")
	io.WriteString(w, "<img width=\"120\" height=\"90\" src=\"https://www.nestle-cereals.com/de/sites/g/files/fawtmp126/files/styles/scale_992/public/d7/packshot_43981575_cini_minis_4x25_green_thread_de_p99_0_2.png?itok=alRje9Ey\"></img>")
}
