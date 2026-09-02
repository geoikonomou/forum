package middleware

import (
	"github.com/geoikonomou/forum/internal/domain/session"
)

type Middleware struct {
	Authorization Authorization
}

func NewMiddleware(sessionManager session.Manager) *Middleware {
	return &Middleware{
		Authorization: NewAuthorizationMiddleware(sessionManager),
	}
}
