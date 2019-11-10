package controller

import (
	"github.com/gorilla/mux"
	"homemadeScreenshotterUploader/controller/actions"
	"homemadeScreenshotterUploader/controller/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.OnPanic)

	router.HandleFunc("/", actions.Default)
	router.HandleFunc("/upload", actions.Upload)
	return router
}
