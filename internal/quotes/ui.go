package quotes

import (
	"net/http"
	"quotes/internal/users"

	"github.com/a-h/templ"
)

type PageContext struct {
	User *users.UserWithPermissions
}

func (s *Service) getPageContext(r *http.Request) PageContext {
	user := s.getUserWithPermissions(r)
	return PageContext{User: user}
}

func (c *PageContext) getUserBadge() templ.Component {
	if c.User == nil {
		return users.UserBadge(nil)
	}
	return users.UserBadge(&c.User.User)
}

func (c *PageContext) getPermissions() users.UserPermissions {
	if c.User == nil {
		return users.UserPermissions{}
	}
	return c.User.Permissions
}
