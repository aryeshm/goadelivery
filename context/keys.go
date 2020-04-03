package context

type key string

const (
	HostKey     = key("hostKey")
	UsernameKey = key("usernameKey")
	PasswordKey = key("passwordKey")

//	databaseKey = key("databaseKey")
)
