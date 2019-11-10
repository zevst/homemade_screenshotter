package actions

import (
	"net/http"
)

// History returns last upload history todo implement
func History(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Default(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
		return
	}

	size := r.FormValue("size")
	if size == "" {
		size = "0"
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(size))
}
