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

func (this *Response) GetRaw() io.Reader {
	return this.response.Body
}

// GetFromJSON response http client
func (this *Response) GetUnmarshalJSON(v interface{}) error {

	err := json.NewDecoder(this.response.Body).Decode(&v)

	if err != nil {
		return err
	}

	return nil
}

// GetFromXML response http client
func (this *Response) GetUnmarshalXML(v interface{}) error {

	err := xml.NewDecoder(this.response.Body).Decode(&v)

	if err != nil {
		return err
	}

	return nil
}

// GetStatusCode http client response
func (this *Response) GetStatusCode() int {
	return this.response.StatusCode
}

// GetHeader http response client
func (this *Response) GetHeader(key string) string {

	return this.response.Header.Get(key)
}

// GetAsString http response client
func (this *Response) GetAsString() (string, error) {
	b, err := ioutil.ReadAll(this.GetRaw())

	if err != nil {
		return "", err
	}

	return string(b), nil
}
