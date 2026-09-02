package middleware

import (
	"testing"

	testhelpers "github.com/geoikonomou/forum/internal/pkg/testing"
)

func TestServices(t *testing.T) {
	mockSessionManager := &testhelpers.MockSessionManager{}

	middleware := NewMiddleware(mockSessionManager)

	auth := middleware.Authorization

	if auth == nil {
		t.Error("Authorization middleware is nil")
	}
}
