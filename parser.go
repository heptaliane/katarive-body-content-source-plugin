package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	pb "github.com/heptaliane/katarive-go-sdk/gen/pb/plugin/v1"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const USER_AGENT string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"

type SourceItemParser struct {
	url string
	doc *goquery.Document
}

func (p *SourceItemParser) Id() string {
	hash := sha256.Sum256([]byte(p.url))
	return hex.EncodeToString(hash[:])
}
func (p *SourceItemParser) Url() string {
	return p.url
}
func (p *SourceItemParser) Title() string {
	return strings.TrimSpace(p.doc.Find("title").Text())
}
func (p *SourceItemParser) Content() string {
	doc := goquery.CloneDocument(p.doc)
	doc.Find("head,script,noscript,style,iframe,header,footer").Remove()
	doc.Find("ruby").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("rt").Text()
		s.ReplaceWithHtml(text)
	})
	return strings.TrimSpace(doc.Text())
}
func NewSourceItemParser(url string) (*SourceItemParser, error) {
	doc, err := GetDocument(url)
	if err != nil {
		return nil, err
	}

	return &SourceItemParser{doc: doc, url: url}, nil
}

type SourceCollectionParser struct {
	parser *SourceItemParser
}

func (p *SourceCollectionParser) Collection() *pb.SourceCollection {
	return nil
}
func (p *SourceCollectionParser) Sources() []*pb.SourceSummary {
	return []*pb.SourceSummary{
		{
			Id:    p.parser.Id(),
			Url:   p.parser.Url(),
			Title: p.parser.Title(),
		},
	}
}
func NewSourceCollectionParser(url string) (*SourceCollectionParser, error) {
	parser, err := NewSourceItemParser(url)
	if err != nil {
		return nil, err
	}
	return &SourceCollectionParser{parser: parser}, nil
}

// helpers
func GetDocument(url string) (*goquery.Document, error) {
	sreq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	sreq.Header.Set("User-Agent", USER_AGENT)

	client := &http.Client{}
	res, err := client.Do(sreq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, &ResponseStatusError{url: url, code: res.StatusCode}
	}

	body := UTF8Body(res)
	return goquery.NewDocumentFromReader(body)
}
func UTF8Body(res *http.Response) io.Reader {
	reader := bufio.NewReader(res.Body)
	head, _ := reader.Peek(1024)
	enc, _, _ := charset.DetermineEncoding(head, res.Header.Get("Content-Type"))
	return transform.NewReader(reader, enc.NewDecoder())
}
