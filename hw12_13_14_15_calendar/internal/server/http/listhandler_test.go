package internalhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	appMocks "github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/server/http/mocks"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

func TestListHandler(t *testing.T) {
	startAt, err := time.Parse(time.RFC3339, "2023-03-12T00:00:00Z")
	require.NoError(t, err)

	events := []storage.Event{{
		ID:          1,
		Title:       "title",
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(time.Hour),
		Description: "description",
		OwnerID:     11,
		SendBefore:  1000,
	}}

	appMock := appMocks.NewApp(t)
	appMock.EXPECT().EventList(mock.Anything, storage.Params{
		StartAtGEq: startAt,
	}).Return(events, nil)
	s := Server{
		app: appMock,
	}

	str := `{
		"startAtGEq": "2023-03-12T00:00:00Z"
	}`
	reader := strings.NewReader(str)
	req, err := http.NewRequest(http.MethodGet, "/api/events", reader) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := s.List()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	res := Response{
		Data: StorageEventListToResponseEventList(events),
	}
	expected, err := json.Marshal(res)
	require.NoError(t, err)
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
