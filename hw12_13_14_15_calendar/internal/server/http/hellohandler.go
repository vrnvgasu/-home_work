package internalhttp

import (
	"net/http"
)

type HelloHandler struct {
	Logger Logger
}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		h.Logger.Error(err.Error())
	}
}
