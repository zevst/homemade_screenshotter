package actions

import "net/http"

// Default shows any std response as "index" action
func Default(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("noaction"))
}
