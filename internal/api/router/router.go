package router

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Call handler constructor here

	// Call mux.handleFunc() here

	return mux
}
