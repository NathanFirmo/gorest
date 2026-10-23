package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Response struct {
	Status     string
	StatusCode int
	Proto      string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
}

var client = &http.Client{Timeout: 60 * time.Second}

func ParseHeaders(raw string) http.Header {
	headers := http.Header{}
	for _, line := range strings.Split(raw, "\n") {
		key, value, ok := strings.Cut(line, ":")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			continue
		}
		headers.Add(key, strings.TrimSpace(value))
	}
	return headers
}

func NormalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url != "" && !strings.Contains(url, "://") {
		return "http://" + url
	}
	return url
}

func MakeRequest(ctx context.Context, method, url, rawHeaders, rawBody string) (*Response, error) {
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(rawBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, NormalizeURL(url), body)
	if err != nil {
		return nil, err
	}
	req.Header = ParseHeaders(rawHeaders)
	if host := req.Header.Get("Host"); host != "" {
		req.Host = host
	}

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Proto:      res.Proto,
		Headers:    res.Header,
		Body:       content,
		Duration:   time.Since(start),
	}, nil
}

func (r *Response) FormattedBody() string {
	if len(r.Body) == 0 {
		return "(empty body)"
	}
	var indented bytes.Buffer
	if json.Indent(&indented, r.Body, "", "  ") == nil {
		return indented.String()
	}
	return string(r.Body)
}

func (r *Response) FormattedHeaders() string {
	keys := make([]string, 0, len(r.Headers))
	for key := range r.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var b strings.Builder
	fmt.Fprintf(&b, "%s %s\n\n", r.Proto, r.Status)
	for _, key := range keys {
		for _, value := range r.Headers[key] {
			fmt.Fprintf(&b, "%s: %s\n", key, value)
		}
	}
	return b.String()
}

func FormatSize(size int) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	}
}

func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2f s", d.Seconds())
}
