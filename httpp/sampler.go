package httpp

import (
	"bytes"
	"harkonnen/errors"
	"harkonnen/runtime"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Sampler struct {
	errorHandler    *runtime.ErrorHandler
	settings        *Settings
	sampleCollector *runtime.SampleCollector
	innerClient     http.Client
	lastResponse    *http.Response
}

func NewSampler(errorHandler *runtime.ErrorHandler, sampleCollector *runtime.SampleCollector, settings *Settings) *Sampler {
	sampler := new(Sampler)
	sampler.errorHandler = errorHandler
	sampler.settings = settings
	sampler.sampleCollector = sampleCollector
	sampler.buildInnerClient()

	return sampler
}

func (sampler *Sampler) UpdateSettings(settings *Settings) {
	sampler.settings = settings
	sampler.buildInnerClient()
}

func (sampler *Sampler) Get(url string, parameters *url.Values) {
	sampler.Request("GET", url, parameters, "", nil)
}

func (sampler *Sampler) Post(url string, contentType string, body io.Reader) {
	sampler.Request("POST", url, nil, contentType, body)
}

func (sampler *Sampler) PostForm(url string, formValues url.Values) {
	sampler.Post(url, "application/x-www-form-urlencoded", strings.NewReader(formValues.Encode()))
}

func (sampler *Sampler) Put(url string, contentType string, body io.Reader) {
	sampler.Request("PUT", url, nil, contentType, body)
}

func (sampler *Sampler) Patch(url string, contentType string, body io.Reader) {
	sampler.Request("PATCH", url, nil, contentType, body)
}

func (sampler *Sampler) Delete(url string, contentType string, body io.Reader) {
	sampler.Request("DELETE", url, nil, contentType, body)
}

func (sampler *Sampler) Request(method string, url string, parameters *url.Values, contentType string, body io.Reader) {
	completeUrl, err := sampler.composeUrl(url)

	if err != nil {
		sampler.errorHandler.Capture(err)
	} else {
		completeUrl := sampler.composeQueryString(completeUrl, parameters)
		request, err := http.NewRequest(method, completeUrl.String(), body)
		sampler.errorHandler.Capture(err)

		if contentType != "" {
			request.Header.Set("Content-Type", contentType)
		}
		sampler.sendRequest(request)
	}
}

func (sampler *Sampler) LastResponse() *http.Response {
	return sampler.lastResponse
}

func (sampler *Sampler) buildInnerClient() {
	client := http.Client{}

	if sampler.settings.KeepCookies {
		client.Jar, _ = cookiejar.New(&cookiejar.Options{})
	}

	if !sampler.settings.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	client.Timeout = sampler.settings.Timeout

	transport := http.Transport{
		TLSHandshakeTimeout:   sampler.settings.TLSHandshakeTimeout,
		DisableKeepAlives:     !sampler.settings.EnableKeepAlive,
		DisableCompression:    !sampler.settings.EnableCompression,
		MaxIdleConns:          sampler.settings.MaxIdleConnections,
		MaxIdleConnsPerHost:   sampler.settings.MaxIdleConnectionsPerHost,
		MaxConnsPerHost:       sampler.settings.MaxConnectionsPerHost,
		IdleConnTimeout:       sampler.settings.IdleConnectionTimeout,
		ResponseHeaderTimeout: sampler.settings.ResponseHeaderTimeout,
	}

	client.Transport = &transport
	sampler.innerClient = client
}

func (sampler *Sampler) composeUrl(url string) (*url.URL, error) {
	if sampler.settings.BaseUrl == "" {
		returnUrl, err := sampler.parseUrl(url)
		return returnUrl, err
	}

	if !strings.HasPrefix(url, "/") {
		sampler.errorHandler.Capture(errors.InvalidPartialUrl{Url: url})
	}

	returnUrl, err := sampler.parseUrl(sampler.settings.BaseUrl + url)
	return returnUrl, err
}

func (sampler *Sampler) parseUrl(rawUrl string) (*url.URL, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		sampler.errorHandler.Capture(errors.MalformedUrl{Url: rawUrl})
	}

	return parsedUrl, err
}

func (sampler *Sampler) composeQueryString(originalAddress *url.URL, params *url.Values) *url.URL {
	if params == nil {
		return originalAddress
	}

	values := originalAddress.Query()

	for key := range *params {
		values.Set(key, params.Get(key))
	}

	originalAddress.RawQuery = values.Encode()
	return originalAddress
}

func (sampler *Sampler) sendRequest(request *http.Request) {
	// Perform the request and track the elapsed time
	startTime := time.Now()
	response, err := sampler.innerClient.Do(request)
	endTime := time.Now()

	sampler.errorHandler.Capture(err)

	// Calculate request and response size
	sentBytes, receivedBytes := sampler.calculateSentReceivedBytes(response)

	// Save query string and strip it from the URL
	pureUrl := request.URL
	queryString := request.URL.Query()
	queryStringForRemoval := request.URL.Query()
	for key := range queryStringForRemoval {
		queryStringForRemoval.Del(key)
	}
	pureUrl.RawQuery = queryStringForRemoval.Encode()

	// Build sample info (specific for HTTP sample)
	info := SampleInfo{
		URL:        pureUrl,
		Parameters: queryString,
		Method:     request.Method,
	}

	// Determine if the request was redirected
	originalURL := request.URL
	finalURL := response.Request.URL
	info.IsRedirect = originalURL != finalURL
	info.FinalURL = finalURL

	// Create request sample
	sample := Sample{
		start:         startTime,
		end:           endTime,
		name:          request.URL.String(),
		sentBytes:     sentBytes,
		receivedBytes: receivedBytes,
		Info:          info,
	}

	sampler.sampleCollector.Collect(sample)
	sampler.lastResponse = response
}

func (sampler *Sampler) calculateSentReceivedBytes(response *http.Response) (int64, int64) {
	// Get original request from response
	request := response.Request

	// Calculate request header size in bytes
	requestHeader, err := httputil.DumpRequestOut(request, false)
	sampler.errorHandler.Capture(err)
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
		sampler.errorHandler.Capture(err)
	}

	// Calculate response header size in bytes
	responseHeader, err := httputil.DumpResponse(response, false)
	sampler.errorHandler.Capture(err)
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
		sampler.errorHandler.Capture(err)
	}

	return requestHeaderSize + requestBodySize, responseHeaderSize + responseBodySize
}
