package internalhttp //nolint:dupl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (s *Server) Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}

		buf := make([]byte, r.ContentLength)
		_, err := r.Body.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			resp.Error.Message = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp)
			return
		}

		req := RequestEvent{}
		if err = json.Unmarshal(buf, &req); err != nil {
			resp.Error.Message = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp)
			return
		}

		event, err := s.app.CreateEvent(context.Background(), req.ToStorageEvent())
		if err != nil {
			resp.Error.Message = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, resp)
			return
		}

		resp.Data = StorageEventToResponseEvent(event)
		w.WriteHeader(http.StatusCreated)
		WriteResponse(w, resp)
	}
}
