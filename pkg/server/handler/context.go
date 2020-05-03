package handler

type ContextKey int

const (
	TokenKey         ContextKey = iota
	AuthorizationKey ContextKey = iota
)
