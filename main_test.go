package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetSourceServiceMetadata(t *testing.T) {
	t.Parallel()

	p := &BodyContentSourceService{Logger: hclog.New(nil)}

	expected := &pb.GetSourceServiceMetadataResponse{
		Name:             "body-content",
		Version:          "v1",
		SupportedPattern: "^https?://.*",
	}

	ctx := context.Background()
	actual, err := p.GetSourceServiceMetadata(ctx, &pb.GetSourceServiceMetadataRequest{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if diff := cmp.Diff(actual, expected, protocmp.Transform()); diff != "" {
		t.Errorf("Unmatched response: %v", diff)
		return
	}
}

func TestGetSource(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(fmt.Sprintf(".%s", r.URL.Path))
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))

	p := &BodyContentSourceService{Logger: hclog.New(nil)}

	cases := map[string]struct {
		path             string
		expectedResponse *pb.GetSourceResponse
		isError          bool
	}{
		"no-title": {
			path: "testdata/no-title.html",
			expectedResponse: &pb.GetSourceResponse{
				Content: "content",
			},
			isError: false,
		},
		"text-only": {
			path: "testdata/text-only.html",
			expectedResponse: &pb.GetSourceResponse{
				Content: "content",
				Title:   "title",
			},
			isError: false,
		},
		"with-script": {
			path: "testdata/with-script.html",
			expectedResponse: &pb.GetSourceResponse{
				Content: "content",
				Title:   "title",
			},
			isError: false,
		},
		"not-html": {
			path: "testdata/not-html.txt",
			expectedResponse: &pb.GetSourceResponse{
				Content: "Not a html",
			},
			isError: false,
		},
		"not-found": {
			path:    "testdata/not-found.html",
			isError: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &pb.GetSourceRequest{
				Url: fmt.Sprintf("%s/%s", ts.URL, tc.path),
			}
			ctx := context.Background()
			data, err := p.GetSource(ctx, req)
			if tc.isError {
				if err == nil {
					t.Error("Error expected but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if diff := cmp.Diff(data, tc.expectedResponse, protocmp.Transform()); diff != "" {
				t.Errorf("Unmatched response: %v", diff)
				return
			}
		})
	}

	t.Cleanup(func() {
		ts.Close()
	})
}
