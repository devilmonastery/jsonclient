package jsonclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestRequest struct {
	In1 string `json:"in_1"`
	In2 string `json:"in_2"`
}

type TestResponse struct {
	Out1 string `json:"out_1"`
	Out2 string `json:"out_2"`
}

func TestCreateClient(t *testing.T) {
	tests := []struct {
		name           string
		request        *TestRequest
		response       *TestResponse
		responseStatus int
		responseBody   string
		expectError    bool
	}{
		{
			name:           "success",
			request:        &TestRequest{In1: "foo", In2: "bar"},
			response:       &TestResponse{Out1: "foo", Out2: "bar"},
			responseStatus: http.StatusOK,
			responseBody:   `{"out_1":"foo","out_2":"bar"}`,
			expectError:    false,
		},
		{
			name:           "server error",
			request:        &TestRequest{In1: "foo", In2: "bar"},
			response:       nil,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `Internal Server Error`,
			expectError:    true,
		},
		{
			name:           "invalid json",
			request:        &TestRequest{In1: "foo", In2: "bar"},
			response:       nil,
			responseStatus: http.StatusOK,
			responseBody:   `not a json`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer svr.Close()

			c := NewJsonClient[TestRequest, TestResponse]()
			x, err := c.Post(svr.URL, tt.request)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("error posting request: %v", err)
			}
			if x.Out1 != tt.response.Out1 || x.Out2 != tt.response.Out2 {
				t.Fatalf("expected response Out1=%s, Out2=%s; got Out1=%s, Out2=%s", tt.response.Out1, tt.response.Out2, x.Out1, x.Out2)
			}
		})
	}
}
