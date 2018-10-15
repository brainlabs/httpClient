package httpClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	DefaultRequestTimeout      = 10000 * time.Millisecond
)

// Client http request very simple
type Client struct {
	Transport *http.Transport
	client    *http.Client
	Cookie    http.CookieJar
	timeout   time.Duration
	headers   map[string]string
	queryUrl  string
}

// NewClient http wrapper
func NewClient() *Client {

	return &Client{
		client: http.DefaultClient,
	}
}

// createClient  handle http client request instance
func (c *Client) createClient() *http.Client {

	if c.client == nil {
		c.client = http.DefaultClient
	}

	if c.Transport == nil {
		// Customize the Transport to have larger connection pool
		defaultTransportPointer, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			log.Fatal("defaultRoundTripper not an *http.Transport")
		}
		defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
		defaultTransport.MaxIdleConns = DefaultMaxIdleConns
		defaultTransport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
		c.client.Transport = &defaultTransport
	}

	c.client.Timeout = c.timeout
	c.client.Jar = c.Cookie

	return c.client
}

// buildUrl  handle http client query url encode
func (c *Client) buildUrl(url string) string {

	if c.queryUrl != "" {
		url = url + c.queryUrl
	}

	return url
}

// SetTimeout handles http client request timeout
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

//// SetCookies handles the receipt of the cookies in a reply for the
//// given URL.  It may or may not choose to save the cookies, depending
//// on the jar's policy and implementation.
//func (c *Client) SetCookies(u *url.URL, cookies []*http.Cookie) {
//	c.lk.Lock()
//	c.cookies[u.Host] = cookies
//	c.lk.Unlock()
//}
//
//// Cookies returns the cookies to send in a request for the given URL.
//// It is up to the implementation to honor the standard cookie use
//// restrictions such as in RFC 6265.
//func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
//	return jar.cookies[u.Host]
//}

// SetHeader  handle http client request header
func (c *Client) SetHeader(key, value string) *Client {

	if c.headers == nil {
		c.headers = map[string]string{
			key: value,
		}

	} else {
		c.headers[key] = value

	}
	return c
}

// SetHeaders  handle http client request headers
func (c *Client) SetHeaders(headers map[string]string) *Client {

	if c.headers == nil {
		c.headers = headers

	} else {
		for k, v := range headers {
			c.headers[k] = v
		}

	}

	return c
}

// SetPemCertificate http client request use ssl
func (c *Client) SetPemCertificate(pemFile string) *Client {

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

	c.Transport = &http.Transport{
		TLSClientConfig:       conf,
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          DefaultMaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		MaxIdleConnsPerHost:   DefaultMaxIdleConnsPerHost,
	}
	return c

}

// SetQuery uri url string  http client request
func (c *Client) SetQuery(query map[string]string) *Client {

	queryUrl := url.Values{}

	for k, v := range query {
		queryUrl.Set(k, v)
	}

	c.queryUrl = "?" + queryUrl.Encode()

	return c
}

// Get method  handle http client request
func (c *Client) Get(url string) (*Response, error) {

	return c.Request("GET", c.buildUrl(url), nil)
}

// Post method  handle http client request
func (c *Client) Post(url string, data []byte) (*Response, error) {

	return c.Request("POST", c.buildUrl(url), data)
}

// Put method  handle http client request
func (c *Client) Put(url string, data []byte) (*Response, error) {

	return c.Request("PUT", c.buildUrl(url), data)
}

// Delete method  handle http client request
func (c *Client) Delete(url string, data []byte) (*Response, error) {

	return c.Request("DELETE", c.buildUrl(url), data)
}

// Request handle http client request
func (c *Client) Request(method, url string, payload []byte) (*Response, error) {
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
	for k, v := range c.headers {
		request.Header.Set(k, v)
	}

	request.Header.Set("Connection", "close")

	if c.timeout < 1 {
		c.timeout = time.Duration(DefaultRequestTimeout)
	}

	// do request client
	client := c.createClient()
	rsp, err := client.Do(request)

	response.response = rsp

	errNet, ok := err.(net.Error);
	if ok && errNet.Timeout() {
		response.isTimeout = ok
		response.statusCode = http.StatusRequestTimeout
	}

	if err != nil {
		return response, err
	}

	return response, nil
}
