package main

import (
	"github.com/daily-demos/daily-prejoin-presence/m/v2/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		retCode  int
		retBody  string
		wantErr  error
		wantRoom Room
	}{
		{
			name:    "success",
			retCode: 200,
			retBody: `
				{
				  "id": "987b5eb5-d116-4a4e-8e2c-14fcb5710966",
				  "name": "presence-test",
				  "api_created": true,
				  "privacy":"private",
				  "url":"https://api-demo.daily.co/presence-test",
				  "created_at":"2019-01-26T09:01:22.000Z",
				  "config":{
					"start_audio_off": true,
					"start_video_off": true
				  }
				}
        `,
			wantRoom: Room{
				Name: "presence-test",
				URL:  "https://api-demo.daily.co/presence-test",
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
			gotRoom, gotErr := createRoom("", testServer.URL)
			require.ErrorIs(t, gotErr, tc.wantErr)
			if tc.wantErr == nil {
				require.EqualValues(t, tc.wantRoom, *gotRoom)
			}
		})

	}
}
