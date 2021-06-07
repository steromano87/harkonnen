package httpp_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"harkonnen/errors"
	"harkonnen/httpp"
	"harkonnen/runtime"
	"harkonnen/telemetry"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type SamplerTestSuite struct {
	suite.Suite
	settings        *httpp.Settings
	errorHandler    *runtime.ErrorHandler
	sampleCollector *runtime.SampleCollector
	sampler         *httpp.Sampler
	testServer      *httptest.Server
}

func (suite *SamplerTestSuite) SetupTest() {
	suite.settings = httpp.NewSettings()
	suite.errorHandler = &runtime.ErrorHandler{}
	suite.sampleCollector = runtime.NewSampleCollector()
	suite.sampler = httpp.NewSampler(suite.errorHandler, suite.sampleCollector, suite.settings)

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		body := string(bodyBytes)

		_, _ = fmt.Fprintf(w, "Request method: '%s'\n", r.Method)
		_, _ = fmt.Fprintf(w, "Request host: 'http://%s'\n", r.Host)
		_, _ = fmt.Fprintf(w, "Request partial URL: '%s'\n", r.URL.String())
		_, _ = fmt.Fprintf(w, "Request body: '%s'\n", body)
	})

	handler.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/redirected")
		w.WriteHeader(302)
	})

	handler.HandleFunc("/redirected", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		body := string(bodyBytes)

		_, _ = fmt.Fprintf(w, "Request method: '%s'\n", r.Method)
		_, _ = fmt.Fprintf(w, "Request host: 'http://%s'\n", r.Host)
		_, _ = fmt.Fprintf(w, "Request partial URL: '%s'\n", r.URL.String())
		_, _ = fmt.Fprintf(w, "Request body: '%s'\n", body)
	})

	suite.testServer = httptest.NewServer(handler)
}

func (suite *SamplerTestSuite) TearDownTest() {
	suite.testServer.Close()
}

func (suite *SamplerTestSuite) TestNewSampler() {
	assert.IsType(suite.T(), &httpp.Sampler{}, suite.sampler)
}

func (suite *SamplerTestSuite) TestGetRequest() {
	suite.sampler.Get(suite.testServer.URL, nil)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "GET", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Equal(suite.T(), url.Values{}, sample.Info.Parameters)
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'GET'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestGetRequestWithQueryString() {
	parameters := url.Values{}
	parameters.Set("key1", "value1")
	parameters.Set("key2", "1")
	suite.sampler.Get(suite.testServer.URL, &parameters)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "GET", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), parameters, sample.Info.Parameters)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'GET'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request partial URL: '/?%s'", parameters.Encode()))
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestPostNoBody() {
	suite.sampler.Post(suite.testServer.URL, "", nil)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "POST", sample.Info.Method)
			// TODO: check if it is better to separate URL from querystring in the Sample
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'POST'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestPutNoBody() {
	suite.sampler.Put(suite.testServer.URL, "", nil)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "PUT", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'PUT'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestPatchNoBody() {
	suite.sampler.Patch(suite.testServer.URL, "", nil)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "PATCH", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'PATCH'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestDeleteNoBody() {
	suite.sampler.Delete(suite.testServer.URL, "", nil)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "DELETE", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'DELETE'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, "Request body: ''")
}

func (suite *SamplerTestSuite) TestPostFormRequest() {
	values := url.Values{}
	values.Set("test", "example")
	suite.sampler.PostForm(suite.testServer.URL, values)
	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		assert.Implements(suite.T(), (*telemetry.Sample)(nil), collectedSamples[0])

		if assert.IsType(suite.T(), httpp.Sample{}, collectedSamples[0]) {
			sample := collectedSamples[0].(httpp.Sample)
			assert.Equal(suite.T(), "POST", sample.Info.Method)
			assert.Equal(suite.T(), suite.testServer.URL, sample.Info.URL.String())
			assert.Equal(suite.T(), suite.testServer.URL, sample.Name())
			assert.Greater(suite.T(), sample.SentBytes(), int64(0))
			assert.Greater(suite.T(), sample.ReceivedBytes(), int64(0))
		}
	}

	assert.IsType(suite.T(), &http.Response{}, suite.sampler.LastResponse())
	responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
	defer func() {
		_ = suite.sampler.LastResponse().Body.Close()
	}()
	responseBody := string(responseBodyBytes)

	assert.Contains(suite.T(), responseBody, "Request method: 'POST'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
	assert.Contains(suite.T(), responseBody, "Request partial URL: '/'")
	assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request body: '%s'", values.Encode()))
}

func (suite *SamplerTestSuite) TestRequestMalformedUrl() {
	malformedUrl := "http:// invalid url"
	suite.sampler.Get(malformedUrl, nil)
	errorList := suite.errorHandler.GetCollected()
	if assert.NotEmpty(suite.T(), errorList) {
		assert.IsType(suite.T(), errors.MalformedUrl{}, errorList[0])
	}
}

func (suite *SamplerTestSuite) TestRequestInvalidPartialUrl() {
	settings := httpp.NewSettings()
	settings.BaseUrl = suite.testServer.URL
	suite.sampler.UpdateSettings(settings)

	invalidPartialUrl := "test"
	suite.sampler.Get(invalidPartialUrl, nil)
	errorList := suite.errorHandler.GetCollected()
	if assert.NotEmpty(suite.T(), errorList) {
		assert.IsType(suite.T(), errors.InvalidPartialUrl{}, errorList[0])
	}
}

func (suite *SamplerTestSuite) TestRequestPartialMalformedUrl() {
	settings := httpp.NewSettings()
	settings.BaseUrl = "https:// my malformed base URL"
	suite.sampler.UpdateSettings(settings)

	malformedPartialUrl := "/test"
	suite.sampler.Get(malformedPartialUrl, nil)
	errorList := suite.errorHandler.GetCollected()
	if assert.NotEmpty(suite.T(), errorList) {
		assert.IsType(suite.T(), errors.MalformedUrl{}, errorList[0])
	}
}

func (suite *SamplerTestSuite) TestRequestWithRedirect_WithoutRedirectSetting() {
	settings := httpp.NewSettings()
	settings.FollowRedirects = false
	settings.BaseUrl = suite.testServer.URL
	suite.sampler.UpdateSettings(settings)

	suite.sampler.Get("/redirect", nil)

	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		sample := collectedSamples[0].(httpp.Sample)
		assert.Equal(suite.T(), suite.testServer.URL+"/redirect", sample.Info.URL.String())

		responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
		defer func() {
			_ = suite.sampler.LastResponse().Body.Close()
		}()
		responseBody := string(responseBodyBytes)

		assert.Empty(suite.T(), responseBody)
	}
}

func (suite *SamplerTestSuite) TestRequestWithRedirect_WithRedirectSetting() {
	settings := httpp.NewSettings()
	settings.FollowRedirects = true
	settings.BaseUrl = suite.testServer.URL
	suite.sampler.UpdateSettings(settings)

	suite.sampler.Get("/redirect", nil)

	assert.Empty(suite.T(), suite.errorHandler.GetCollected())

	collectedSamples := suite.sampleCollector.Flush()

	if assert.Equal(suite.T(), 1, len(collectedSamples)) {
		sample := collectedSamples[0].(httpp.Sample)
		assert.Equal(suite.T(), suite.testServer.URL+"/redirect", sample.Info.URL.String())
		assert.True(suite.T(), sample.Info.IsRedirect)
		assert.Equal(suite.T(), suite.testServer.URL+"/redirected", sample.Info.FinalURL.String())

		responseBodyBytes, _ := ioutil.ReadAll(suite.sampler.LastResponse().Body)
		defer func() {
			_ = suite.sampler.LastResponse().Body.Close()
		}()
		responseBody := string(responseBodyBytes)

		assert.Contains(suite.T(), responseBody, "Request method: 'GET'")
		assert.Contains(suite.T(), responseBody, fmt.Sprintf("Request host: '%s'", suite.testServer.URL))
		assert.Contains(suite.T(), responseBody, "Request partial URL: '/redirected'")
		assert.Contains(suite.T(), responseBody, "Request body: ''")
	}
}

func TestSamplerTestSuite(t *testing.T) {
	suite.Run(t, new(SamplerTestSuite))
}
