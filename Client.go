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
func (c *Client) createClient() *http.Client {
	if c.client == nil {
		c.client = &http.Client{
			Transport: &c.Transport,
			Jar:       c.Cookie,
			Timeout:   c.timeout,
		}
	}

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

	c.Transport = http.Transport{
		TLSClientConfig:       conf,
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
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

	if c.timeout < 1 {
		c.timeout = time.Duration(10 * time.Second)
	}

	// do request client
	client := c.createClient()
	rsp, err := client.Do(request)

	if err != nil {
		return response, err
	}

	response.response = rsp

	return response, nil
}
