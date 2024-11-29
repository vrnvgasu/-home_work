package internalhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return IPAddress
}

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type RequestEvent struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	Description string    `json:"description"`
	OwnerID     uint64    `json:"ownerId"`
	SendBefore  int64     `json:"sendBefore"`
}

type ResponseEvent struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	Description string    `json:"description"`
	OwnerID     uint64    `json:"ownerId"`
	SendBefore  int64     `json:"sendBefore"`
}

type RequestEventList struct {
	Events []RequestEvent `json:"events"`
}

func (r *RequestEvent) ToStorageEvent() storage.Event {
	return storage.Event{
		ID:          r.ID,
		Title:       r.Title,
		StartAt:     r.StartAt,
		EndAt:       r.EndAt,
		Description: r.Description,
		OwnerID:     r.OwnerID,
		SendBefore:  r.SendBefore,
	}
}

func StorageEventToResponseEvent(e storage.Event) ResponseEvent {
	return ResponseEvent{
		ID:          e.ID,
		Title:       e.Title,
		StartAt:     e.StartAt,
		EndAt:       e.EndAt,
		Description: e.Description,
		OwnerID:     e.OwnerID,
		SendBefore:  e.SendBefore,
	}
}

func StorageEventListToResponseEventList(evs []storage.Event) []ResponseEvent {
	res := make([]ResponseEvent, 0, len(evs))
	for _, e := range evs {
		res = append(res, StorageEventToResponseEvent(e))
	}

	return res
}

type ParamsRequest struct {
	Limit      int
	StartAtGEq time.Time `json:"startAtGEq"`
	StartAtLEq time.Time `json:"startAtLEq"`
}

func (p *ParamsRequest) ToStorageParams() storage.Params {
	return storage.Params{
		Limit:      p.Limit,
		StartAtGEq: p.StartAtGEq,
		StartAtLEq: p.StartAtLEq,
	}
}

func WriteResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
