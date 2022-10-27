package msnp

type msnp_context struct {
	dispatched bool
	ctxkey     int
}

var msn_context_list []*msnp_context
