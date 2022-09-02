package cservice

type (
	ContextKey string
)

const (
	ContextKeyMetaKey      ContextKey = "X-MetaKey"
	ContextKeyProcessID    ContextKey = "X-ProcessID"
	ContextKeyProcessName  ContextKey = "X-ProcessName"
	ContextKeyServiceStack ContextKey = "X-ServiceStack"
)
