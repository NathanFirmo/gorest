package curl

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    Request
	}{
		{
			name:    "simple get",
			command: "curl https://api.example.com/users",
			want:    Request{Method: "GET", URL: "https://api.example.com/users"},
		},
		{
			name:    "explicit method and headers",
			command: `curl -X PUT 'https://api.example.com/users/1' -H 'Authorization: Bearer abc' -H "Accept: application/json"`,
			want: Request{
				Method: "PUT",
				URL:    "https://api.example.com/users/1",
				Headers: []Header{
					{Key: "Authorization", Value: "Bearer abc"},
					{Key: "Accept", Value: "application/json"},
				},
			},
		},
		{
			name:    "attached short flag values",
			command: `curl -XDELETE -H'X-Id: 7' https://api.example.com/items/7`,
			want: Request{
				Method:  "DELETE",
				URL:     "https://api.example.com/items/7",
				Headers: []Header{{Key: "X-Id", Value: "7"}},
			},
		},
		{
			name: "browser copy as cURL with line continuations",
			command: "curl 'https://api.example.com/login' \\\n" +
				"  -H 'content-type: application/json' \\\n" +
				"  --data-raw $'{\"user\":\"nathan\",\"note\":\"it\\'s ok\"}' \\\n" +
				"  --compressed",
			want: Request{
				Method:  "POST",
				URL:     "https://api.example.com/login",
				Headers: []Header{{Key: "content-type", Value: "application/json"}},
				Body:    `{"user":"nathan","note":"it's ok"}`,
			},
		},
		{
			name:    "form data defaults content type",
			command: `curl https://example.com/form -d a=1 -d b=2`,
			want: Request{
				Method:  "POST",
				URL:     "https://example.com/form",
				Headers: []Header{{Key: "Content-Type", Value: "application/x-www-form-urlencoded"}},
				Body:    "a=1&b=2",
			},
		},
		{
			name:    "json flag",
			command: `curl --json '{"a":1}' https://example.com`,
			want: Request{
				Method: "POST",
				URL:    "https://example.com",
				Headers: []Header{
					{Key: "Content-Type", Value: "application/json"},
					{Key: "Accept", Value: "application/json"},
				},
				Body: `{"a":1}`,
			},
		},
		{
			name:    "basic auth and user agent",
			command: `curl -u user:pass -A gorest -sSL https://example.com`,
			want: Request{
				Method: "GET",
				URL:    "https://example.com",
				Headers: []Header{
					{Key: "Authorization", Value: "Basic dXNlcjpwYXNz"},
					{Key: "User-Agent", Value: "gorest"},
				},
			},
		},
		{
			name:    "get moves data to query string",
			command: `curl -G https://example.com/search?lang=go --data-urlencode 'q=hello world'`,
			want:    Request{Method: "GET", URL: "https://example.com/search?lang=go&q=hello+world"},
		},
		{
			name:    "head request ignores output flags",
			command: `curl -I -o /dev/null --max-time 5 https://example.com`,
			want:    Request{Method: "HEAD", URL: "https://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.command)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, command := range []string{
		"wget https://example.com",
		"curl -H 'Accept: */*'",
		"curl 'https://example.com",
		"curl -F file=@a.png https://example.com",
		"curl https://example.com -X",
	} {
		if _, err := Parse(command); err == nil {
			t.Errorf("Parse(%q) expected an error", command)
		}
	}
}

func TestHeaderText(t *testing.T) {
	req := Request{Headers: []Header{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}}
	if got := req.HeaderText(); got != "A: 1\nB: 2" {
		t.Errorf("HeaderText() = %q", got)
	}
}
