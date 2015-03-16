package token

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fishworks/api/auth"
)

func TestIsAuthenticated(t *testing.T) {
	token := New(&auth.User{})
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("HTTP_AUTHORIZATION", fmt.Sprintf("token %s", token.Key))
	if !token.IsAuthenticated(req) {
		t.Error("expected req to be authenticated")
	}
	req.Header.Set("HTTP_AUTHORIZATION", "token bad")
	if token.IsAuthenticated(req) {
		t.Error("expected req to not be authenticated")
	}
}
