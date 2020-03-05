package domadoma

import "net/http"

func renderError(w http.ResponseWriter,msg string,statuscode int) {
	w.WriteHeader(statuscode)
	w.Write([]byte(msg))
}