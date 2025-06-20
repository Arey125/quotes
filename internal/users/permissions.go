package users

type Permisson string

const (
	PermissonQuotesRead       Permisson = "quotes.read"
	PermissonQuotesWrite      Permisson = "quotes.write"
	PermissonQuotesModeration Permisson = "quotes.moderation"
	PermissonUserPermissions  Permisson = "users.permissions"
)
