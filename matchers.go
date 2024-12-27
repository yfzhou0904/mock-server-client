package mockserver

type RequestMatcher struct {
	Method                string              `json:"method,omitempty"`
	Path                  string              `json:"path,omitempty"`
	Body                  *BodyMatcher        `json:"body,omitempty"`
	Headers               map[string][]string `json:"headers,omitempty"`
	QueryStringParameters map[string][]string `json:"queryStringParameters,omitempty"`
}

func (m RequestMatcher) WithHeader(key, value string) RequestMatcher {
	if m.Headers == nil {
		m.Headers = make(map[string][]string)
	}
	m.Headers[key] = []string{value}
	return m
}

type BodyMatcher struct {
	// complete spec can be found at https://www.mock-server.com/mock_server/creating_expectations.html

	// JSON | PARAMETERS | STRING | REGEX | BINARY
	Type string `json:"type,omitempty"`

	// when type is "JSON"
	MatchType string                 `json:"matchType,omitempty"`
	JSON      map[string]interface{} `json:"json,omitempty"`

	// when type is "PARAMETERS"
	// aka form fields, body parameters, or "application/x-www-form-urlencoded"
	Parameters map[string][]string `json:"parameters,omitempty"`

	// when type is "STRING"
	// strict string equality check on request body
	String string `json:"string,omitempty"`

	// when type is "REGEX"
	// regex match on request body
	Regex string `json:"regex,omitempty"`

	// when type is "BINARY"
	Base64Bytes string `json:"base64Bytes,omitempty"`
}

const (
	MatchBodyJSON = "JSON"
)
const (
	StrictMatch   = "STRICT"
	TolerantMatch = "ONLY_MATCHING_FIELDS"
)

func (m RequestMatcher) WithJSONFields(json map[string]interface{}) RequestMatcher {
	m.Body = &BodyMatcher{
		Type:      MatchBodyJSON,
		JSON:      json,
		MatchType: TolerantMatch,
	}
	return m
}
