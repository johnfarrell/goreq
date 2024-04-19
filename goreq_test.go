package goreq

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestParse(t *testing.T) {
	type SimpleRequest struct {
		QueryParam string
	}
	type args struct {
		Queries map[string]string
	}
	tests := []struct {
		name     string
		args     args
		want     SimpleRequest
		wantCode int
	}{
		{
			name: "Simple request",
			args: args{
				Queries: map[string]string{
					"queryparam": "value",
				},
			},
			want: SimpleRequest{
				QueryParam: "value",
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "Bad request",
			args:     args{},
			want:     SimpleRequest{},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, ok := GetRequest[SimpleRequest](r.Context())
			if !ok {
				t.Fatalf("Request not found in request context")
			}

			if req.QueryParam != tt.want.QueryParam {
				t.Errorf("RequestParse() got=%v, want %v", req.QueryParam, tt.want.QueryParam)
			}
		})

		testHandler := RequestParse[SimpleRequest](nextHandler)

		req := buildRequest(tt.args.Queries, nil)

		rr := httptest.NewRecorder()
		testHandler.ServeHTTP(rr, req)

		if rr.Code != tt.wantCode {
			t.Errorf("RequestParse() got=%v, want %v", rr.Code, tt.wantCode)
		}
	}
}
