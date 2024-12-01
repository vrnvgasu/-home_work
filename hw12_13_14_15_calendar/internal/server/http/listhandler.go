package internalhttp //nolint:dupl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (s *Server) List() http.HandlerFunc {
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

		req := ParamsRequest{}
		if err = json.Unmarshal(buf, &req); err != nil {
			resp.Error.Message = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			WriteResponse(w, resp)
			return
		}

		events, err := s.app.EventList(context.Background(), req.ToStorageParams())
		if err != nil {
			resp.Error.Message = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, resp)
			return
		}

		resp.Data = StorageEventListToResponseEventList(events)
		w.WriteHeader(http.StatusOK)
		WriteResponse(w, resp)
	}
}
