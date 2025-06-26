package config

import (
	"testing"
	"titan-lift/internal/server_error"
)

func TestBadFile(t *testing.T) {
	t.Run("test bad file", func(t *testing.T) {
		data := []byte("bad data")
		_, err := GetServerConfig(data)

		if !server_error.IsServerError(err, "CONFIG_PARSER") {
			t.Error("Failing bad file")
		}
	})
}
