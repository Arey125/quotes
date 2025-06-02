package users

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type contextKey string

var contextKeyUser = contextKey("user")

type InjectUserMiddleware struct {
	model          *Model
	sessionManager *scs.SessionManager
}

func NewInjectUserMiddleware(model *Model, sessionManager *scs.SessionManager) InjectUserMiddleware {
	return InjectUserMiddleware{model: model, sessionManager: sessionManager}
}

func (m *InjectUserMiddleware) Wrap(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userWithPermissions := func() *UserWithPermissions {
			userId, ok := m.sessionManager.Get(r.Context(), "user_id").(int)
			if !ok {
				return nil
			}
			user, err := m.model.GetUserWithPermissions(userId)
			if err != nil || user == nil {
				return nil
			}
			return user
		}()

		ctx := context.WithValue(r.Context(), contextKeyUser, userWithPermissions)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetUser(r *http.Request) *UserWithPermissions {
	user, ok := r.Context().Value(contextKeyUser).(*UserWithPermissions)
	if !ok {
		panic("could not get user from request context, probably InjectUserMiddleware is absent")
	}
	return user
}
