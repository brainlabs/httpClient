package httpClient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	response   *http.Response
	isTimeout  bool
	statusCode int
}

var noReader = &emptyReader{}

// emptyReader
type emptyReader struct {
}

func (ep *emptyReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("the reponse of request  is nil")
}

func (r *Response) GetRaw() io.Reader {

	if r.response == nil {
		return noReader
	}

	bd := r.response.Body

	r.response.Body.Close()

	return  bd
}

// GetFromJSON response http client
func (r *Response) GetUnmarshalJSON(v interface{}) error {

	if r.GetStatusCode() == http.StatusRequestTimeout {
		return http.ErrHandlerTimeout
	}

	err := json.NewDecoder(r.GetRaw()).Decode(&v)

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

	err := xml.NewDecoder(r.GetRaw()).Decode(&v)

	if err != nil {
		return err
	}

	return nil
}

// GetStatusCode http client response
func (r *Response) GetStatusCode() int {

	if r.response != nil {
		r.statusCode = r.response.StatusCode
	}
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

	b, err := ioutil.ReadAll(r.GetRaw())

	if err != nil {
		return result, err
	}

	result = string(b)

	return result, nil
}
