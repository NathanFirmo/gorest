package utils

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseHeaders(t *testing.T) {
	headers := ParseHeaders("Authorization: Bearer a:b\n\ninvalid\n Referer : https://example.com \n: empty")
	if got := headers.Get("Authorization"); got != "Bearer a:b" {
		t.Errorf("Authorization = %q", got)
	}
	if got := headers.Get("Referer"); got != "https://example.com" {
		t.Errorf("Referer = %q", got)
	}
	if len(headers) != 2 {
		t.Errorf("expected 2 headers, got %v", headers)
	}
}

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Method", r.Method)
		w.Header().Set("X-Token", r.Header.Get("X-Token"))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"echo":` + string(body) + `}`))
	}))
	defer server.Close()

	res, err := MakeRequest(context.Background(), "POST", strings.TrimPrefix(server.URL, "http://"), "X-Token: abc", `{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusCreated || res.Headers.Get("X-Method") != "POST" || res.Headers.Get("X-Token") != "abc" {
		t.Fatalf("unexpected response %+v", res)
	}
	if want := "{\n  \"echo\": {\n    \"a\": 1\n  }\n}"; res.FormattedBody() != want {
		t.Errorf("FormattedBody() = %q, want %q", res.FormattedBody(), want)
	}
}

func TestFormattedBodyKeepsNonJSON(t *testing.T) {
	res := &Response{Body: []byte("<h1>hi</h1>")}
	if res.FormattedBody() != "<h1>hi</h1>" {
		t.Errorf("FormattedBody() = %q", res.FormattedBody())
	}
}
