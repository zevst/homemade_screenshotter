package middleware

import (
	"fmt"
	"net/http"
	"os"
)

// OnPanic is Mux's middleware to handle panic situations
func OnPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer onActionPanic(w)

		next.ServeHTTP(w, r)
	})
}

// panic handler for controller actions
func onActionPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		var message string
		switch x := r.(type) {
		case string:
			message = x
		case error:
			message = x.Error()
		default:
			message = fmt.Sprintf("unknown error '%s'", x)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(message))
		_, _ = fmt.Fprint(os.Stderr, message)
	}
}
