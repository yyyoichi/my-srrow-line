package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

var key = []byte(os.Getenv("CSRF_TOKEN"))
var CSRFMiddleware = csrf.Protect(key,
	csrf.SameSite(csrf.SameSiteStrictMode),
	csrf.Secure(true),
	csrf.HttpOnly(true),
	csrf.ErrorHandler(http.HandlerFunc(csrfError)),
	csrf.CookieName("_csrf_token"),
	csrf.MaxAge(2),
)

func csrfError(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("invalid token from '%s' method '%s' to '%s'", r.Host, r.Method, r.RequestURI)
	http.Error(w, "error", http.StatusForbidden)
}
