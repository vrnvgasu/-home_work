//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/stretchr/testify/require"
)

type getRequest struct {
	StartAtGEq time.Time `json:"startAtGEq"`
}

type getResponse struct {
	Data  []getEntity `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type getEntity struct {
	ID          uint64    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	StartAt     time.Time `json:"startAt" db:"start_at"`
	EndAt       time.Time `json:"endAt" db:"end_at"`
	Description string    `json:"description" db:"description"`
	OwnerID     uint64    `json:"ownerId" db:"owner_id"`
	SendBefore  int64     `json:"sendBefore" db:"send_before"`
}

func (s *SuiteIntegrationTest) TestGet() {
	_, err := s.db.Exec(`insert into event (title, start_at, end_at, description, owner_id, send_before) 
		values ('title_day', now() - interval '10 hours' , now(), '', '1', '100'),
		 ('title_week', now() - interval '6 days' , now(), '', '1', '100'),
		 ('title_month', now() - interval '15 days' , now(), '', '1', '100') `)
	s.Require().NoError(err)

	cases := []struct {
		name               string
		startAtGEq         time.Time
		itemsCount         int
		expectedTitleNames []string
	}{
		{
			name:               "event for day",
			startAtGEq:         time.Now().Add(-24 * time.Hour),
			itemsCount:         1,
			expectedTitleNames: []string{"title_day"},
		},
		{
			name:               "event for week",
			startAtGEq:         time.Now().Add(-7 * 24 * time.Hour),
			itemsCount:         2,
			expectedTitleNames: []string{"title_day", "title_week"},
		},
		{
			name:               "event for month",
			startAtGEq:         time.Now().Add(-30 * 24 * time.Hour),
			itemsCount:         3,
			expectedTitleNames: []string{"title_day", "title_week", "title_month"},
		},
	}

	for _, tt := range cases {
		tt := tt
		s.Run(tt.name, func() {
			response := s.httpRequest(http.MethodGet, "/api/events", getRequest{
				StartAtGEq: tt.startAtGEq,
			})
			defer response.Body.Close()

			getBody, err := io.ReadAll(response.Body)
			s.Require().NoError(err)

			var resp getResponse
			err = json.Unmarshal(getBody, &resp)
			s.Require().NoError(err)

			require.Len(s.T(), resp.Data, tt.itemsCount)

			titleNames := make([]string, 0, tt.itemsCount)
			for _, item := range resp.Data {
				titleNames = append(titleNames, item.Title)
			}

			sort.Strings(titleNames)
			sort.Strings(tt.expectedTitleNames)
			require.Equal(s.T(), tt.expectedTitleNames, titleNames)
		})
	}
}
