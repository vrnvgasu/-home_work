//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/stretchr/testify/require"
)

type createResponse struct {
	Data  eventEntity `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type eventEntity struct {
	ID          uint64    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	StartAt     time.Time `json:"startAt" db:"start_at"`
	EndAt       time.Time `json:"endAt" db:"end_at"`
	Description string    `json:"description" db:"description"`
	OwnerID     uint64    `json:"ownerId" db:"owner_id"`
	SendBefore  int64     `json:"sendBefore" db:"send_before"`
}

type createRequest struct {
	Title       string    `json:"title"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	Description string    `json:"description"`
	OwnerID     uint64    `json:"ownerId"`
	SendBefore  int64     `json:"sendBefore"`
}

func (s *SuiteIntegrationTest) TestAdd() {
	startAt, err := time.Parse(time.RFC3339, "2023-03-12T00:00:00Z")
	s.Require().NoError(err)
	endAt, err := time.Parse(time.RFC3339, "2023-04-12T00:00:00Z")
	s.Require().NoError(err)

	eventRequest := createRequest{
		Title:       "http_test",
		StartAt:     startAt,
		EndAt:       endAt,
		Description: "description",
		OwnerID:     1,
		SendBefore:  10000,
	}
	response := s.httpRequest(http.MethodPost, "/api/events", eventRequest)
	defer response.Body.Close()

	getBody, err := io.ReadAll(response.Body)
	s.Require().NoError(err)

	var resp createResponse
	err = json.Unmarshal(getBody, &resp)
	s.Require().NoError(err)

	require.Equal(s.T(), eventRequest.Title, resp.Data.Title)
	require.Equal(s.T(), eventRequest.StartAt, resp.Data.StartAt)
	require.Equal(s.T(), eventRequest.EndAt, resp.Data.EndAt)
	require.Equal(s.T(), eventRequest.Description, resp.Data.Description)
	require.Equal(s.T(), eventRequest.OwnerID, resp.Data.OwnerID)
	require.Equal(s.T(), eventRequest.SendBefore, resp.Data.SendBefore)

	if !s.needCheckSql() {
		return
	}

	var eventDB eventEntity
	rows, err := s.db.NamedQuery("select id, title, start_at, end_at, description, owner_id, send_before from event where id = :id", map[string]any{
		"id": resp.Data.ID,
	})
	s.Require().NoError(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&eventDB)
		s.Require().NoError(err)
	}

	eventDB.StartAt.Format(time.RFC3339)
	eventDB.EndAt.Format(time.RFC3339)

	require.Equal(s.T(), eventDB.ID, resp.Data.ID)
	require.Equal(s.T(), eventDB.Title, resp.Data.Title)
	require.True(s.T(), eventDB.StartAt.Equal(resp.Data.StartAt))
	require.True(s.T(), eventDB.EndAt.Equal(resp.Data.EndAt))
	require.Equal(s.T(), eventDB.Description, resp.Data.Description)
	require.Equal(s.T(), eventDB.OwnerID, resp.Data.OwnerID)
	require.Equal(s.T(), eventDB.SendBefore, resp.Data.SendBefore)
}
