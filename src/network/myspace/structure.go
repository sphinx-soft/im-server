package myspace

type MySpaceDataPair struct {
	Key   string
	Value string
}

type MySpaceContext struct {
	Nonce      string
	SessionKey int
	Status     MySpaceStatus
}

type MySpaceStatus struct {
	Code    int
	Message string
}

type MySpaceUserDetails struct {
	UIN           int
	FavouriteBand string
	FavouriteSong string
	Age           int
	Gender        string
	Location      string
}

var clientContexts []*MySpaceContext
