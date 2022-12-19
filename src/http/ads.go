package http

import (
	"io"
	"math/rand"
	"net/http"
	"phantom/util"
)

func CycleMySpaceAds(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	urls := []string{"http://lu.is-very-gay.lol/content/cdn/25Px4pz8iXnk.png", "http://www.nestle-cereals.com/de/sites/g/files/fawtmp126/files/styles/scale_992/public/d7/packshot_43981575_cini_minis_4x25_green_thread_de_p99_0_2.png"}
	randomIndex := rand.Intn(len(urls))
	pick := urls[randomIndex]
	util.Log(util.INFO, "AdServer", "Sending Advertisment Data!")
	io.WriteString(w, "<img width=\"120\" height=\"90\" src=\""+pick+"\"></img>")
}
