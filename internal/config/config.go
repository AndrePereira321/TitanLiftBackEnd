package config

import (
	"bytes"
	"github.com/spf13/viper"
	"titan-lift/internal/server_error"
)

type ServerConfig struct {
	httpServer HttpConfig
}

func (c *ServerConfig) HttpServer() HttpConfig {
	return c.httpServer
}

type HttpConfig struct {
	host    string
	port    uint
	sslPort uint
}

func (s *HttpConfig) Host() string {
	return s.host
}

func (s *HttpConfig) Port() uint {
	return s.port
}

func GetServerConfig(data []byte) (*ServerConfig, error) {
	v := getViper()
	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, server_error.Wrap("CONFIG_PARSER", "error when reading config data", err)
	}

	serverConfig := &ServerConfig{}
	if httpConfig, err := getHttpConfig(v); err != nil {
		return nil, err
	} else {
		serverConfig.httpServer = *httpConfig
	}

	return serverConfig, nil
}

func getHttpConfig(v *viper.Viper) (*HttpConfig, error) {
	httpConfig := HttpConfig{}

	httpConfig.host = v.GetString("server.host")
	if len(httpConfig.host) == 0 {
		return nil, server_error.New("CONFIG_PARSER", "http server host is empty")
	}

	httpConfig.port = v.GetUint("server.port")
	if httpConfig.port == 0 {
		httpConfig.port = 80
	}

	httpConfig.sslPort = v.GetUint("server.ssl_port")
	if httpConfig.sslPort == 0 {
		httpConfig.sslPort = 443
	}

	return &httpConfig, nil
}

func getViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("toml")
	return v
}
