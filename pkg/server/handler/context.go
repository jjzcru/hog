package handler

// ContextKey unique identifier
type ContextKey int

const (
	// TokenKey is a key for a token
	TokenKey ContextKey = iota
	// AuthorizationKey is a key for an authorization header
	AuthorizationKey ContextKey = iota
)
