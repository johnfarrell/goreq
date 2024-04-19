package goreq

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func buildRequest(queries, headers map[string]string) *http.Request {
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
func Test_parseTags(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want structTags
	}{
		{
			name: "empty",
			tag:  "",
			want: structTags{
				In: inValQuery,
			},
		},
		{
			name: "empty value",
			tag:  "in=",
			want: structTags{
				In: inValQuery,
			},
		},
		{
			name: "invalid param",
			tag:  "in",
			want: structTags{
				In: inValQuery,
			},
		},
		{
			name: "empty and valid value",
			tag:  "in=query,label=",
			want: structTags{
				In: inValQuery,
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
				In:    inValQuery,
				Label: "value",
			},
		},
		{
			name: "Simple type",
			tag:  "type=json",
			want: structTags{
				In:   inValQuery,
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
				t.Errorf("parseTags() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_parseParameters_BasicRequestWithTags(t *testing.T) {
	type SampleRequest struct {
		StringHeader string `goreq:"in=header,label=string-header"`
		StringQuery  string `goreq:"in=query,label=string-query"`
		SkippedValue string `goreq:"-"` // Would be processed as a query param if not skipped
	}

	type args struct {
		Queries map[string]string
		Headers map[string]string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *T
		wantErr bool
	}
	tests := []testCase[SampleRequest]{
		{
			name: "Empty request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Header-only request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{
					"string-header": "value",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Query-only request",
			args: args{
				Queries: map[string]string{
					"string-query": "value",
				},
				Headers: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Complete request",
			args: args{
				Queries: map[string]string{
					"string-query": "value",
				},
				Headers: map[string]string{
					"string-header": "value",
				},
			},
			want: &SampleRequest{
				StringHeader: "value",
				StringQuery:  "value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseParameters[SampleRequest](buildRequest(tt.args.Queries, tt.args.Headers))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParameters() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseParameters() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
func Test_parseParameters_BasicRequestNoTags(t *testing.T) {
	// With no tags, these would both be expected in the query params
	type SampleRequest struct {
		Header string
		Query  string
	}

	type args struct {
		Queries map[string]string
		Headers map[string]string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *T
		wantErr bool
	}
	tests := []testCase[SampleRequest]{
		{
			name: "Empty request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Header-only request",
			args: args{
				Queries: map[string]string{},
				Headers: map[string]string{
					"Header": "value",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "One query specified",
			args: args{
				Queries: map[string]string{
					"query": "value",
				},
				Headers: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Complete request",
			args: args{
				Queries: map[string]string{
					"query":  "value",
					"header": "foo",
				},
				Headers: map[string]string{},
			},
			want: &SampleRequest{
				Query:  "value",
				Header: "foo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseParameters[SampleRequest](buildRequest(tt.args.Queries, tt.args.Headers))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParameters() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseParameters() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
