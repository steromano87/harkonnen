package rest

import (
	"bytes"
	"github.com/steromano87/harkonnen/shooter"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"time"
)

type Client struct {
	context      shooter.Context
	settings     *Settings
	innerClient  http.Client
	lastResponse *http.Response
}

func NewClient(context shooter.Context, settings *Settings) *Client {
	client := new(Client)
	client.context = context
	client.settings = settings
	client.buildInnerClient()

	return client
}

func (c *Client) UpdateSettings(settings *Settings) {
	c.settings = settings
	c.buildInnerClient()
}

func (c *Client) LastResponse() *http.Response {
	return c.lastResponse
}

func (c *Client) Execute(request Request, options ...Option) {
	// Generate the raw request
	rawRequest, err := request.Build(c.settings.BaseUrl)

	if err != nil {
		c.context.OnUnrecoverableError(err)
		return
	}

	// Perform the request and track the elapsed time
	startTime := time.Now()
	response, err := c.innerClient.Do(rawRequest)
	endTime := time.Now()

	if err != nil {
		c.context.OnUnrecoverableError(err)
		return
	}

	// Calculate request and response size
	sentBytes, receivedBytes := c.calculateSentReceivedBytes(response)

	// Save query string and strip it from the URL
	pureUrl := rawRequest.URL
	queryString := rawRequest.URL.Query()
	queryStringForRemoval := rawRequest.URL.Query()
	for key := range queryStringForRemoval {
		queryStringForRemoval.Del(key)
	}

	pureUrl.RawQuery = queryStringForRemoval.Encode()
	originalURL := rawRequest.URL
	finalURL := response.Request.URL

	// Create request sample
	sample := NewSample(rawRequest.URL.String(), startTime, endTime, sentBytes, receivedBytes)
	sample.URL = pureUrl
	sample.Parameters = queryString
	sample.Method = request.Method
	sample.IsRedirect = originalURL != finalURL
	sample.FinalURL = finalURL

	c.context.SampleCollector().Collect(sample)
	c.lastResponse = response
}

func (c *Client) calculateSentReceivedBytes(response *http.Response) (int64, int64) {
	// Get original request from response
	request := response.Request

	// Calculate request header size in bytes
	requestHeader, err := httputil.DumpRequestOut(request, false)
	if err != nil {
		c.context.OnUnrecoverableError(err)
	}

	requestHeaderSize := int64(len(requestHeader))

	// Calculate request body size in bytes
	// Using a double buffer to prevent the body to be consumed by the byte count operation
	// See https://stackoverflow.com/a/23077519
	requestBodySize := int64(0)
	if request.Body != nil {
		bodyBuffer, err := ioutil.ReadAll(request.Body)
		countWriter := ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		backupWriter := ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		requestBodySize, err = io.Copy(io.Discard, countWriter)
		request.Body = backupWriter

		if err != nil {
			c.context.OnUnrecoverableError(err)
		}
	}

	// Calculate response header size in bytes
	responseHeader, err := httputil.DumpResponse(response, false)
	if err != nil {
		c.context.OnUnrecoverableError(err)
	}

	responseHeaderSize := int64(len(responseHeader))

	// Calculate response body size in bytes
	// Using a double buffer to prevent the body to be consumed by the byte count operation
	// See https://stackoverflow.com/a/23077519
	responseBodySize := int64(0)
	if response.Body != nil {
		bodyBuffer, err := ioutil.ReadAll(response.Body)
		countWriter := ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		backupWriter := ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
		responseBodySize, err = io.Copy(io.Discard, countWriter)
		response.Body = backupWriter

		if err != nil {
			c.context.OnUnrecoverableError(err)
		}
	}

	return requestHeaderSize + requestBodySize, responseHeaderSize + responseBodySize
}

func (c *Client) buildInnerClient() {
	client := http.Client{}

	if c.settings.KeepCookies {
		client.Jar, _ = cookiejar.New(&cookiejar.Options{})
	}

	if !c.settings.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	client.Timeout = c.settings.Timeout

	transport := http.Transport{
		TLSHandshakeTimeout:   c.settings.TLSHandshakeTimeout,
		DisableKeepAlives:     !c.settings.EnableKeepAlive,
		DisableCompression:    !c.settings.EnableCompression,
		MaxIdleConns:          c.settings.MaxIdleConnections,
		MaxIdleConnsPerHost:   c.settings.MaxIdleConnectionsPerHost,
		MaxConnsPerHost:       c.settings.MaxConnectionsPerHost,
		IdleConnTimeout:       c.settings.IdleConnectionTimeout,
		ResponseHeaderTimeout: c.settings.ResponseHeaderTimeout,
	}

	client.Transport = &transport
	c.innerClient = client
}
