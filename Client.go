package httpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client http request very simple
type Client struct {
	Transport http.Transport
	client    *http.Client
	Cookie    http.CookieJar
	timeout   time.Duration
	headers   map[string]string
	queryUrl  string
}

func NewClient() *Client {

	return &Client{}
}

// createClient  handle http client request instance
func (this *Client) createClient() *http.Client {
	if this.client == nil {
		this.client = &http.Client{
			Transport: &this.Transport,
			Jar:       this.Cookie,
			Timeout:   this.timeout,
		}
	}

	return this.client
}

// buildUrl  handle http client query url encode
func (this *Client) buildUrl(url string) string {

	if this.queryUrl != "" {
		url = url + this.queryUrl
	}

	return url
}

// SetTimeout handles http client request timeout
func (this *Client) SetTimeout(timeout time.Duration) *Client {
	this.timeout = timeout
	return this
}

//// SetCookies handles the receipt of the cookies in a reply for the
//// given URL.  It may or may not choose to save the cookies, depending
//// on the jar's policy and implementation.
//func (this *Client) SetCookies(u *url.URL, cookies []*http.Cookie) {
//	this.lk.Lock()
//	this.cookies[u.Host] = cookies
//	this.lk.Unlock()
//}
//
//// Cookies returns the cookies to send in a request for the given URL.
//// It is up to the implementation to honor the standard cookie use
//// restrictions such as in RFC 6265.
//func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
//	return jar.cookies[u.Host]
//}

// SetHeader  handle http client request header
func (this *Client) SetHeader(key, value string) *Client {

	if this.headers == nil {
		this.headers = map[string]string{
			key: value,
		}

	} else {
		this.headers[key] = value

	}
	return this
}

// SetHeaders  handle http client request headers
func (this *Client) SetHeaders(headers map[string]string) *Client {

	if this.headers == nil {
		this.headers = headers

	} else {
		for k, v := range headers {
			this.headers[k] = v
		}

	}

	return this
}

// SetPemCertificate http client request use ssl
func (this *Client) SetPemCertificate(pemFile string) *Client {

	cert, err := ioutil.ReadFile(pemFile)
	if err != nil {
		log.Fatalf("couldn't load file pem ", err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	conf := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: false,
	}

	defaultTransport := http.DefaultTransport.(*http.Transport)

	this.Transport = http.Transport{
		TLSClientConfig:       conf,
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
	}
	return this

}

// SetQuery uri url string  http client request
func (this *Client) SetQuery(query map[string]string) *Client {

	queryUrl := url.Values{}

	for k, v := range query {
		queryUrl.Set(k, v)
	}

	this.queryUrl = "?" + queryUrl.Encode()

	return this
}

// Get method  handle http client request
func (this *Client) Get(url string) (*Response, error) {

	return this.Request("GET", this.buildUrl(url), nil)
}

// Post method  handle http client request
func (this *Client) Post(url string, data []byte) (*Response, error) {

	return this.Request("POST", this.buildUrl(url), data)
}

// Put method  handle http client request
func (this *Client) Put(url string, data []byte) (*Response, error) {

	return this.Request("PUT", this.buildUrl(url), data)
}

// Delete method  handle http client request
func (this *Client) Delete(url string, data []byte) (*Response, error) {

	return this.Request("DELETE", this.buildUrl(url), data)
}

// Request handle http client request
func (this *Client) Request(method, url string, payload []byte) (*Response, error) {
	var request *http.Request
	var err error

	response := &Response{}

	method = strings.ToUpper(method)
	request, err = http.NewRequest(method, url, bytes.NewBuffer(payload))

	// if make request error
	if err != nil {
		return response, err
	}

	// create header
	for k, v := range this.headers {
		request.Header.Set(k, v)
	}

	if this.timeout < 1 {
		this.timeout = time.Duration(10 * time.Second)
	}

	// do request client
	client := this.createClient()
	rsp, err := client.Do(request)

	if err != nil {
		return response, err
	}

	response.response = rsp

	return response, nil
}
