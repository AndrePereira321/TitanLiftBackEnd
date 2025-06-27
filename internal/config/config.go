package config

import (
	"bytes"
	"github.com/spf13/viper"
	"titan-lift/internal/server_error"
)

type ServerConfig struct {
	httpServer HttpConfig
	logging    LoggingConfig
}

func (c ServerConfig) HttpServer() HttpConfig {
	return c.httpServer
}

func (c ServerConfig) Logging() LoggingConfig {
	return c.logging
}

type HttpConfig struct {
	host    string
	port    uint
	sslPort uint
}

func (s HttpConfig) Host() string {
	return s.host
}

func (s HttpConfig) Port() uint {
	return s.port
}

type LoggingConfig struct {
	logDir         string
	serverLogLevel string
	queryLogLevel  string
}

func (l LoggingConfig) LogDir() string {
	return l.logDir
}

func (l LoggingConfig) ServerLogLevel() string {
	return l.serverLogLevel
}

func (l LoggingConfig) QueryLogLevel() string {
	return l.queryLogLevel
}

func GetServerConfig(data []byte) (*ServerConfig, error) {
	serverConfig := &ServerConfig{}
	v := getViper()

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, server_error.Wrap("CONFIG_PARSER", "error when reading config data", err)
	}

	if httpConfig, err := getHttpConfig(v); err != nil {
		return nil, err
	} else {
		serverConfig.httpServer = *httpConfig
	}

	if loggingConfig, err := getLoggingConfig(v); err != nil {
		return nil, err
	} else {
		serverConfig.logging = *loggingConfig
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

func getLoggingConfig(v *viper.Viper) (*LoggingConfig, error) {
	loggingConfig := LoggingConfig{}

	loggingConfig.logDir = v.GetString("logging.log_dir")

	loggingConfig.serverLogLevel = v.GetString("logging.server_log_level")
	if len(loggingConfig.serverLogLevel) == 0 {
		loggingConfig.serverLogLevel = "INFO"
	}

	loggingConfig.queryLogLevel = v.GetString("logging.query_log_level")
	if len(loggingConfig.queryLogLevel) == 0 {
		loggingConfig.queryLogLevel = "INFO"
	}

	return &loggingConfig, nil
}

func getViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("toml")
	return v
}
