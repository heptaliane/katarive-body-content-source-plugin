package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	katarive "github.com/heptaliane/katarive-go-sdk"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
)

const NAME string = "body-content"
const VERSION string = "v1"
const SUPPORTED_PATTERN string = "^https?://.*"

type BodyContentSourceService struct {
	pb.UnimplementedSourceServiceServer
	logger hclog.Logger
}

func (s *BodyContentSourceService) GetSourceServiceMetadata(
	ctx context.Context,
	req *pb.GetSourceServiceMetadataRequest,
) (*pb.GetSourceServiceMetadataResponse, error) {
	s.logger.Trace("GetSourceServiceMetadata called")
	return &pb.GetSourceServiceMetadataResponse{
		Name:             NAME,
		Version:          VERSION,
		SupportedPattern: SUPPORTED_PATTERN,
	}, nil
}
func (s *BodyContentSourceService) GetSource(
	ctx context.Context,
	req *pb.GetSourceRequest,
) (*pb.GetSourceResponse, error) {
	s.logger.Trace("GetSource called", "url", req.GetUrl())
	res, err := http.Get(req.GetUrl())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &ResponseStatusError{url: req.GetUrl(), code: res.StatusCode}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	doc.Find("head,script,noscript,style,iframe,header,footer").Remove()
	content := strings.TrimSpace(doc.Text())

	return &pb.GetSourceResponse{
		Title:   title,
		Content: content,
	}, nil
}

// Ensure BodyContentSourceService implements SourceServiceServer
var _ pb.SourceServiceServer = new(BodyContentSourceService)

type ResponseStatusError struct {
	url  string
	code int
}

func (e *ResponseStatusError) Error() string {
	return fmt.Sprintf("Status code error for '%s': %d", e.url, e.code)
}

// Ensure ResponseStatusError implements error
var _ error = new(ResponseStatusError)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Trace,
		Output: os.Stderr,
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: katarive.Handshake,
		Plugins: map[string]plugin.Plugin{
			"source": &katarive.SourcePlugin{
				Impl: &BodyContentSourceService{
					logger: logger,
				},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})
}
