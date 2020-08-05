package mockserver

type Expectation struct {
	Request  RequestMatcher `json:"httpRequest,omitempty"`
	Response Response       `json:"httpResponse,omitempty"`
	Priority int32          `json:"priority,omitempty"`
}

func NewExpectation(request RequestMatcher) *Expectation {
	e := new(Expectation)
	e.Request = request
	return e
}

func (e Expectation) WithResponse(response Response) Expectation {
	e.Response = response
	return e
}

func (e Expectation) WithPriority(priority int32) Expectation {
	e.Priority = priority
	return e
}

type Response struct {
	Body       map[string]interface{} `json:"body,omitempty"`
	Headers    map[string][]string    `json:"headers,omitempty"`
	StatusCode int32                  `json:"statusCode,omitempty"`
}

func NewResponseOK() *Response {
	e := new(Response)
	e.StatusCode = 200
	return e
}

func (e Response) WithJSONBody(json map[string]interface{}) Response {
	e.Body = json
	return e
}

func (e Response) WithHeader(key, value string) Response {
	if e.Headers == nil {
		e.Headers = make(map[string][]string)
	}
	e.Headers[key] = []string{value}
	return e
}
