package graphql

type ContextKey int

const (
	TokenKey         ContextKey = iota
	AuthorizationKey ContextKey = iota
)
