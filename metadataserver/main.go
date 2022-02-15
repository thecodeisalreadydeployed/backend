package main

import "net/http"

type metadataHandler struct{}

func (handler *metadataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &metadataHandler{})
	http.ListenAndServe(":5000", mux)
}
