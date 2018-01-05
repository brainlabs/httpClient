package httpClient

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	response *http.Response
}

func (r *Response) GetRaw() io.Reader {
	return r.response.Body
}

// GetFromJSON response http client
func (r *Response) GetUnmarshalJSON(v interface{}) error {

	if r.GetStatusCode() == http.StatusRequestTimeout {
		return http.ErrHandlerTimeout
	}
	err := json.NewDecoder(r.response.Body).Decode(&v)

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
	err := xml.NewDecoder(r.response.Body).Decode(&v)

	if err != nil {
		return err
	}

	return nil
}

// GetStatusCode http client response
func (r *Response) GetStatusCode() int {
	return r.response.StatusCode
}

// GetHeader http response client
func (r *Response) GetHeader(key string) string {

	return r.response.Header.Get(key)
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
