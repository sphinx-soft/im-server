package http

import "net/http"

//"net/http"

func HandleAPI(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		/*
			global := util.GetClientFromGlobalList(r.URL.Query().Get("username"))

			if global != nil {

				if r.URL.Query().Has("client") {
					io.WriteString(w, global.Client+" | Protocol: "+global.Protocol)
					return
				}

				if r.URL.Query().Has("friends") {
					io.WriteString(w, strconv.Itoa(global.Friends))
					return
				}
			}
		*/
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
