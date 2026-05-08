package main

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const USER_AGENT string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36"

type SourceParser struct {
	doc *goquery.Document
}

func (p *SourceParser) Title() string {
	return strings.TrimSpace(p.doc.Find("title").Text())
}
func (p *SourceParser) Content() string {
	doc := goquery.CloneDocument(p.doc)
	doc.Find("head,script,noscript,style,iframe,header,footer").Remove()
	doc.Find("ruby").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("rt").Text()
		s.ReplaceWithHtml(text)
	})
	return strings.TrimSpace(doc.Text())
}
func NewSourceParser(url string) (*SourceParser, error) {
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
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	return &SourceParser{doc: doc}, nil
}

// helpers
func UTF8Body(res *http.Response) io.Reader {
	reader := bufio.NewReader(res.Body)
	head, _ := reader.Peek(1024)
	enc, _, _ := charset.DetermineEncoding(head, res.Header.Get("Content-Type"))
	return transform.NewReader(reader, enc.NewDecoder())
}
