package main_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	plugin "github.com/heptaliane/katarive-body-content-source-plugin"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetSourceServiceMetadata(t *testing.T) {
	t.Parallel()

	p := new(plugin.BodyContentSourceService)

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
