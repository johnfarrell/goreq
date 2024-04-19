package goreq

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_parseTags(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want structTags
	}{
		{
			name: "empty",
			tag:  "",
			want: structTags{},
		},
		{
			name: "empty value",
			tag:  "in=",
			want: structTags{},
		},
		{
			name: "invalid param",
			tag:  "in",
			want: structTags{},
		},
		{
			name: "empty and valid value",
			tag:  "in=query,label=",
			want: structTags{
				In: "query",
			},
		},
		{
			name: "Simple in",
			tag:  "in=value",
			want: structTags{
				In: "value",
			},
		},
		{
			name: "Simple label",
			tag:  "label=value",
			want: structTags{
				Label: "value",
			},
		},
		{
			name: "Simple type",
			tag:  "type=json",
			want: structTags{
				Type: "json",
			},
		},
		{
			name: "multiple values",
			tag:  "in=header,label=x-user-email",
			want: structTags{
				In:    "header",
				Label: "x-user-email",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTags(tt.tag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseParameters_BasicRequest(t *testing.T) {
	buildRequest := func(queries, headers map[string]string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

		q := req.URL.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		return req
	}

	type SampleRequest struct {
		StringHeader string `goreq:"in=header,label=string-header"`
		StringQuery  string `goreq:"in=query,label=string-query"`
	}

	type args struct {
		Queries map[string]string
		Headers map[string]string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[SampleRequest]{
		{
			name: "Empty request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{},
			},
			want:    SampleRequest{},
			wantErr: false,
		},
		{
			name: "Header-only request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{
					"string-header": "value",
				},
			},
			want: SampleRequest{
				StringHeader: "value",
			},
			wantErr: false,
		},
		{
			name: "Query-only request",
			args: args{
				Queries: map[string]string{
					"string-query": "value",
				},
				Headers: map[string]string{},
			},
			want: SampleRequest{
				StringQuery: "value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseParameters[SampleRequest](buildRequest(tt.args.Queries, tt.args.Headers))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseParameters() got = %v, want %v", got, tt.want)
			}
		})
	}
}
