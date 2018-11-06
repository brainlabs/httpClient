package httpClient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	response   *http.Response
	isTimeout  bool
	statusCode int
	byteBody   []byte
	err        error
}

var noReader = &emptyReader{}

// emptyReader
type emptyReader struct {
}

func NewResponse(rsp *http.Response) *Response {
	x := &Response{
		response: rsp,
	}
	x.initialize()
	return x
}

func (r *Response) initialize() {
	if r.response != nil {
		r.statusCode = r.response.StatusCode
		r.byteBody, r.err = ioutil.ReadAll(r.GetRaw())
		defer r.response.Body.Close()
	}
}

func (ep *emptyReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("the reponse of request is nil")
}

func (r *Response) GetRaw() []byte {

	return r.byteBody
}

// GetFromJSON response http client
func (r *Response) GetUnmarshalJSON(v interface{}) error {

	if r.GetStatusCode() == http.StatusRequestTimeout {
		return http.ErrHandlerTimeout
	}
	if r.err != nil {
		return r.err
	}

	err := json.Unmarshal(r.byteBody, &v)

	if err != nil {
		return err
	}

	return nil
}

// GetFromXML response http client
func (r *Response) GetUnmarshalXML(v interface{}) error {

	if r.GetStatusCode() == http.StatusRequestTimeout {
		return http.ErrHandlerTimeout
	}

	if r.err != nil {
		return r.err
	}

	err := xml.Unmarshal(r.byteBody, &v)

	if err != nil {
		return err
	}

	return nil
}

// GetStatusCode http client response
func (r *Response) GetStatusCode() int {
	return r.statusCode
}

// GetHeader http response client
func (r *Response) GetHeader(key string) string {

	if r.GetStatusCode() == 0 {
		return ""
	}

	return r.response.Header.Get(key)
}

// IsTimeout request response
func (r *Response) IsTimeout() bool {
	return r.isTimeout
}

// GetAsString http response client
func (r *Response) GetAsString() (string, error) {
	var result string

	if r.GetStatusCode() == http.StatusRequestTimeout {
		return result, http.ErrHandlerTimeout
	}

	result = string(r.byteBody)

	return result, r.err
}
