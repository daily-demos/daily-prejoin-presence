package main

import (
	"github.com/daily-demos/daily-prejoin-presence/m/v2/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPresence(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		retCode          int
		retBody          string
		wantErr          error
		wantParticipants []Participant
	}{
		{
			name:    "one-participant",
			retCode: 200,
			// This response is copied directly from the presence endpoint docs example:
			// https://docs.daily.co/reference/rest-api/rooms/get-room-presence#example-request
			retBody: `
				{
				  "total_count": 1,
				  "data": [
					{
					  "room": "w2pp2cf4kltgFACPKXmX",
					  "id": "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					  "userId": "pbZ+ismP7dk=",
					  "userName": "Moishe",
					  "joinTime": "2023-01-01T20:53:19.000Z",
					  "duration": 2312
					}
				  ]
				}
        	`,
			wantParticipants: []Participant{
				{
					ID:   "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					Name: "Moishe",
				},
			},
		},
		{
			name:    "three-participants",
			retCode: 200,
			// This response is copied directly from the presence endpoint docs example:
			// https://docs.daily.co/reference/rest-api/rooms/get-room-presence#example-request
			retBody: `
				{
				  "total_count": 3,
				  "data": [
					{
					  "room": "w2pp2cf4kltgFACPKXmX",
					  "id": "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					  "userId": "pbZ+ismP7dk=",
					  "userName": "Moishe",
					  "joinTime": "2023-01-01T20:53:19.000Z",
					  "duration": 2312
					},
					{
					  "room": "w2pp2cf4kltgFACPKXmX",
					  "id": "participant-id",
					  "userId": "participant-id",
					  "userName": "Liza",
					  "joinTime": "2023-01-01T20:53:19.000Z",
					  "duration": 2312
					},
					{
					  "room": "w2pp2cf4kltgFACPKXmX",
					  "id": "participant-id-2",
					  "userId": "participant-id-2",
					  "userName": "Bob",
					  "joinTime": "2023-01-01T20:53:19.000Z",
					  "duration": 2312
					}
				  ]
				}
        	`,
			wantParticipants: []Participant{
				{
					ID:   "d61cd7b2-a273-42b4-89bd-be763fd562c1",
					Name: "Moishe",
				},
				{
					ID:   "participant-id",
					Name: "Liza",
				},
				{
					ID:   "participant-id-2",
					Name: "Bob",
				},
			},
		},
		{
			name:    "failure",
			retCode: 500,
			wantErr: util.ErrFailedDailyAPICall,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.retCode)
				_, err := w.Write([]byte(tc.retBody))
				require.NoError(t, err)
			}))
			defer testServer.Close()
			gotParticipants, gotErr := getPresence("name", "key", testServer.URL)
			require.ErrorIs(t, gotErr, tc.wantErr)
			if tc.wantErr == nil {
				require.EqualValues(t, tc.wantParticipants, gotParticipants)
			}
		})

	}
}
