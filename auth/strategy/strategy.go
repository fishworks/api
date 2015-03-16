package strategy

import (
	"net/http"
)

// Strategy defines how HTTP requests are authenticated against different auth backends.
type Strategy interface {
	IsAuthenticated(req *http.Request) bool
}
