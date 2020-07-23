package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/qubole/wireguard/internal/testhelpers"
)

func Test_Socks_Server(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Health Check",
			url:        "/socks_healthz",
			wantStatus: http.StatusOK,
			wantBody:   `{"name":"socks proxy","status":"OK"}`,
		},
		{
			name:       "Test Failure",
			url:        "/wontWork",
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   "",
		},
	}

	socks := NewSocks(Port(8090))
	socks.logger = testhelpers.FakeLogger(false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("--- %v ", err.Error())
			}
			rr := httptest.NewRecorder()
			socks.router.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("Wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			if tt.wantBody != "" {
				if got := rr.Body.String(); got != tt.wantBody {
					t.Errorf("Wrong response body: got %v want %v",
						got, tt.wantBody)
				}
			}
		})
	}
}
