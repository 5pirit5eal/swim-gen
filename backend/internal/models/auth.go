package models

type ctxKey string

const UserIdCtxKey ctxKey = "user_id"

type SharingMethod string

const (
	SharingMethodLink  SharingMethod = "link"
	SharingMethodEmail SharingMethod = "email"
)
