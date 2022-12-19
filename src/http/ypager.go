package http

import (
	"net/http"
	"phantom/util"
)

/*
 // cookie will get expired after 1 year
    expires := time.Now().AddDate(1, 0, 0)

    ck := http.Cookie{
        Name: "JSESSION_ID",
        Domain: "foo.com",
        Path: "/",
        Expires: expires,
    }

    // value of cookie
    ck.Value = "value of this awesome cookie"

    // write the cookie to response
    http.SetCookie(w, &ck)
*/

func HandleYPager(w http.ResponseWriter, r *http.Request) {
	util.Log(util.INFO, "Yahoo! Pager", "Client awaiting Web Authentication from: %s", r.RemoteAddr)

	//cookie := http.Cookie{
	//		Name: "Y",
	//		Value: "",
	//	}
}
