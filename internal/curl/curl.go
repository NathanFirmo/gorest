package curl

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode"
)

type Header struct {
	Key   string
	Value string
}

type Request struct {
	Method  string
	URL     string
	Headers []Header
	Body    string
}

var ErrNotCurl = errors.New("not a curl command")

var valueFlags = map[string]string{
	"-X": "request", "--request": "request",
	"-H": "header", "--header": "header",
	"-d": "data", "--data": "data", "--data-raw": "data", "--data-binary": "data", "--data-ascii": "data",
	"--data-urlencode": "data-urlencode",
	"--json":           "json",
	"-u":               "user", "--user": "user",
	"-A": "user-agent", "--user-agent": "user-agent",
	"-b": "cookie", "--cookie": "cookie",
	"-e": "referer", "--referer": "referer",
	"--url": "url",
	"-F":    "form", "--form": "form",

	"-o": "", "--output": "",
	"-m": "", "--max-time": "", "--connect-timeout": "",
	"-w": "", "--write-out": "",
	"-x": "", "--proxy": "",
	"-c": "", "--cookie-jar": "",
	"-E": "", "--cert": "", "--key": "", "--cacert": "",
	"-K": "", "--config": "",
	"-r": "", "--range": "",
	"-T": "", "--upload-file": "",
	"--retry": "", "--resolve": "", "--limit-rate": "",
}

func (r Request) HeaderText() string {
	lines := make([]string, len(r.Headers))
	for i, h := range r.Headers {
		lines[i] = h.Key + ": " + h.Value
	}
	return strings.Join(lines, "\n")
}

func (r Request) header(key string) (string, bool) {
	for _, h := range r.Headers {
		if strings.EqualFold(h.Key, key) {
			return h.Value, true
		}
	}
	return "", false
}

func (r *Request) setDefaultHeader(key, value string) {
	if _, ok := r.header(key); !ok {
		r.Headers = append(r.Headers, Header{Key: key, Value: value})
	}
}

func IsCurl(command string) bool {
	fields := strings.Fields(command)
	return len(fields) > 0 && fields[0] == "curl"
}

func Parse(command string) (Request, error) {
	args, err := tokenize(command)
	if err != nil {
		return Request{}, err
	}
	if len(args) == 0 || args[0] != "curl" {
		return Request{}, ErrNotCurl
	}

	var (
		req           Request
		data          []string
		head, get     bool
		isJSON        bool
		explicitBlank bool
	)

	for i := 1; i < len(args); i++ {
		arg := args[i]
		flag, value, hasValue := splitFlag(arg)

		if !strings.HasPrefix(arg, "-") || arg == "-" {
			if req.URL == "" {
				req.URL = arg
			}
			continue
		}

		name, takesValue := valueFlags[flag]
		if !takesValue {
			switch arg {
			case "-I", "--head":
				head = true
			case "-G", "--get":
				get = true
			}
			continue
		}

		if !hasValue {
			if i+1 >= len(args) {
				return Request{}, fmt.Errorf("missing value for %s", flag)
			}
			i++
			value = args[i]
		}

		switch name {
		case "request":
			req.Method = strings.ToUpper(value)
		case "header":
			key, val, ok := strings.Cut(value, ":")
			if ok && strings.TrimSpace(key) != "" {
				req.Headers = append(req.Headers, Header{Key: strings.TrimSpace(key), Value: strings.TrimSpace(val)})
			}
		case "data":
			data = append(data, value)
			explicitBlank = explicitBlank || value == ""
		case "data-urlencode":
			data = append(data, urlEncodeData(value))
		case "json":
			data = append(data, value)
			isJSON = true
		case "user":
			user, pass, _ := strings.Cut(value, ":")
			token := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
			req.Headers = append(req.Headers, Header{Key: "Authorization", Value: "Basic " + token})
		case "user-agent":
			req.Headers = append(req.Headers, Header{Key: "User-Agent", Value: value})
		case "cookie":
			if strings.Contains(value, "=") {
				req.Headers = append(req.Headers, Header{Key: "Cookie", Value: value})
			}
		case "referer":
			req.Headers = append(req.Headers, Header{Key: "Referer", Value: value})
		case "url":
			req.URL = value
		case "form":
			return Request{}, errors.New("multipart forms (-F) are not supported yet")
		}
	}

	if req.URL == "" {
		return Request{}, errors.New("curl command has no URL")
	}

	hasData := len(data) > 0 || explicitBlank
	body := strings.Join(data, "&")

	if get && hasData {
		separator := "?"
		if strings.Contains(req.URL, "?") {
			separator = "&"
		}
		req.URL += separator + body
		hasData = false
		body = ""
	}

	if req.Method == "" {
		switch {
		case head:
			req.Method = "HEAD"
		case hasData:
			req.Method = "POST"
		default:
			req.Method = "GET"
		}
	}

	if hasData {
		req.Body = body
		if isJSON {
			req.setDefaultHeader("Content-Type", "application/json")
			req.setDefaultHeader("Accept", "application/json")
		} else {
			req.setDefaultHeader("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	return req, nil
}

// Short flags accept an attached value, e.g. -XPOST or -H'Accept: */*'.
func splitFlag(arg string) (flag, value string, hasValue bool) {
	if strings.HasPrefix(arg, "--") || len(arg) <= 2 || arg[0] != '-' {
		return arg, "", false
	}
	if _, ok := valueFlags[arg[:2]]; ok {
		return arg[:2], arg[2:], true
	}
	return arg, "", false
}

func urlEncodeData(value string) string {
	name, content, ok := strings.Cut(value, "=")
	if !ok {
		return url.QueryEscape(value)
	}
	if name == "" {
		return url.QueryEscape(content)
	}
	return name + "=" + url.QueryEscape(content)
}

// tokenize splits a command line following POSIX shell quoting rules, plus
// bash's $'...' strings, which browsers emit when copying requests as cURL.
func tokenize(command string) ([]string, error) {
	var (
		args    []string
		current strings.Builder
		inToken bool
	)
	runes := []rune(command)

	for i := 0; i < len(runes); i++ {
		c := runes[i]
		switch {
		case c == '\\':
			if i+1 >= len(runes) {
				continue
			}
			i++
			if runes[i] == '\r' && i+1 < len(runes) && runes[i+1] == '\n' {
				i++
				continue
			}
			if runes[i] == '\n' {
				continue
			}
			current.WriteRune(runes[i])
			inToken = true
		case c == '\'':
			end := indexRune(runes, i+1, '\'')
			if end < 0 {
				return nil, errors.New("unterminated single quote")
			}
			current.WriteString(string(runes[i+1 : end]))
			i = end
			inToken = true
		case c == '$' && i+1 < len(runes) && runes[i+1] == '\'':
			end, err := readANSIC(runes, i+2, &current)
			if err != nil {
				return nil, err
			}
			i = end
			inToken = true
		case c == '"':
			i++
			for ; i < len(runes) && runes[i] != '"'; i++ {
				if runes[i] == '\\' && i+1 < len(runes) && strings.ContainsRune("\"\\$`\n", runes[i+1]) {
					i++
					if runes[i] == '\n' {
						continue
					}
				}
				current.WriteRune(runes[i])
			}
			if i >= len(runes) {
				return nil, errors.New("unterminated double quote")
			}
			inToken = true
		case unicode.IsSpace(c):
			if inToken {
				args = append(args, current.String())
				current.Reset()
				inToken = false
			}
		default:
			current.WriteRune(c)
			inToken = true
		}
	}

	if inToken {
		args = append(args, current.String())
	}
	return args, nil
}

func readANSIC(runes []rune, start int, out *strings.Builder) (int, error) {
	escapes := map[rune]rune{'n': '\n', 't': '\t', 'r': '\r', '\\': '\\', '\'': '\'', '"': '"', 'e': '\x1b', '0': 0}
	for i := start; i < len(runes); i++ {
		switch runes[i] {
		case '\'':
			return i, nil
		case '\\':
			if i+1 >= len(runes) {
				continue
			}
			i++
			if r, ok := escapes[runes[i]]; ok {
				out.WriteRune(r)
			} else {
				out.WriteRune('\\')
				out.WriteRune(runes[i])
			}
		default:
			out.WriteRune(runes[i])
		}
	}
	return 0, errors.New("unterminated $'...' string")
}

func indexRune(runes []rune, from int, target rune) int {
	for i := from; i < len(runes); i++ {
		if runes[i] == target {
			return i
		}
	}
	return -1
}
