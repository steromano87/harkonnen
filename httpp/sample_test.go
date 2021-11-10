package httpp_test

import (
	"github.com/stretchr/testify/assert"
	"harkonnen/httpp"
	"net/url"
	"testing"
	"time"
)

var testUrl, _ = url.Parse("https://www.wikipedia.org")

var testHttpSample = httpp.NewSample(
	"Test HTTP Sample",
	time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
	32,
	2048,
	httpp.SampleInfo{
		URL:    testUrl,
		Method: "GET",
	},
)

func TestSample_Name(t *testing.T) {
	assert.Equal(t, "Test HTTP Sample", testHttpSample.Name())
}

func TestSample_Start(t *testing.T) {
	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), testHttpSample.Start())
}

func TestSample_End(t *testing.T) {
	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), testHttpSample.End())
}

func TestSample_Duration(t *testing.T) {
	assert.Equal(t, 1*time.Second, testHttpSample.Duration())
}

func TestSample_SentBytes(t *testing.T) {
	assert.EqualValues(t, 32, testHttpSample.SentBytes())
}

func TestSample_ReceivedBytes(t *testing.T) {
	assert.EqualValues(t, 2048, testHttpSample.ReceivedBytes())
}

func TestSampleInfo(t *testing.T) {
	info := testHttpSample.Info

	assert.IsType(t, httpp.SampleInfo{}, info)
	assert.Equal(t, testUrl, info.URL)
	assert.Equal(t, "GET", info.Method)
}