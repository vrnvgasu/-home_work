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

func TestUpdateHandler(t *testing.T) {
	startAt, err := time.Parse(time.RFC3339, "2023-03-12T00:00:00Z")
	require.NoError(t, err)
	endAt, err := time.Parse(time.RFC3339, "2023-04-12T00:00:00Z")
	require.NoError(t, err)

	event := storage.Event{
		ID:          1,
		Title:       "title",
		StartAt:     startAt,
		EndAt:       endAt,
		Description: "description",
		OwnerID:     11,
		SendBefore:  1000,
	}

	appMock := appMocks.NewApp(t)
	appMock.EXPECT().UpdateEvent(mock.Anything, event).Return(nil)
	s := Server{
		app: appMock,
	}

	str := `{
		  "id": 1,
		  "title": "title",
		  "startAt": "2023-03-12T00:00:00Z",
		  "endAt": "2023-04-12T00:00:00Z",
		  "description": "description",
		  "ownerId": 11,
		  "sendBefore": 1000
	}`
	reader := strings.NewReader(str)
	req, err := http.NewRequest(http.MethodPut, "/api/events", reader) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := s.Update()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	res := Response{
		Data: StorageEventToResponseEvent(event),
	}
	expected, err := json.Marshal(res)
	require.NoError(t, err)
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
