package msim

type msim_data_pair struct {
	Key   string
	Value string
}

type msim_context struct {
	nonce         string
	sesskey       int
	statuscode    int
	statusmessage string
	details       msim_user_details
}

type msim_user_details struct {
	userid     int
	avatartype string
	bandname   string
	songname   string
	age        int
	gender     string
	location   string
	headline   string
	lastlogin  int64
}

var users_context []*msim_context
