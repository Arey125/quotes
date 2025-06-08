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

func (c *PageContext) getNavigation() templ.Component {
	return users.Navigation(c.User)
}

func (c *PageContext) getPermissions() users.UserPermissions {
	if c.User == nil {
		return users.UserPermissions{}
	}
	return c.User.Permissions
}
