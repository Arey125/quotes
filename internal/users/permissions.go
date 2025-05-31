package users

type Permisson string

const (
	PermissonQuotesRead      Permisson = "quotes.read"
	PermissonQuotesWrite     Permisson = "quotes.write"
	PermissonUserPermissions Permisson = "users.permissions"
)
