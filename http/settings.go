package http

import (
	"time"
)

type Settings struct {
	Timeout                   time.Duration `mapstructure:"timeout"`
	BaseUrl                   string        `mapstructure:"baseUrl"`
	FollowRedirects           bool          `mapstructure:"followRedirects"`
	KeepCookies               bool          `mapstructure:"keepCookies"`
	IdleConnectionTimeout     time.Duration `mapstructure:"idleConnectionTimeout"`
	TLSHandshakeTimeout       time.Duration `mapstructure:"tlsHandshakeTimeout"`
	ResponseHeaderTimeout     time.Duration `mapstructure:"responseHeaderTimeout"`
	MaxIdleConnections        int           `mapstructure:"maxIdleConnections"`
	MaxConnectionsPerHost     int           `mapstructure:"maxConnectionPerHost"`
	MaxIdleConnectionsPerHost int           `mapstructure:"maxIdleConnectionsPerHost"`
	EnableKeepAlive           bool          `mapstructure:"enableKeepAlive"`
	EnableCompression         bool          `mapstructure:"enableCompression"`
}

func NewSettings() *Settings {
	settings := new(Settings)
	settings.Timeout = 30 * time.Second
	settings.BaseUrl = ""
	settings.FollowRedirects = true
	settings.KeepCookies = true
	settings.IdleConnectionTimeout = 10 * time.Second
	settings.TLSHandshakeTimeout = 10 * time.Second
	settings.ResponseHeaderTimeout = 10 * time.Second
	settings.MaxIdleConnections = 100
	settings.MaxConnectionsPerHost = 100
	settings.MaxIdleConnectionsPerHost = 100
	settings.EnableKeepAlive = true
	settings.EnableCompression = true

	return settings
}
