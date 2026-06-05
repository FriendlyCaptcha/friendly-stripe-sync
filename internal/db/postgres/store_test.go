package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHostList(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		isValid bool
	}{
		{"single host", "localhost", true},
		{"comma-separated", "10.0.0.1,10.0.0.2", true},
		{"whitespace", "10.0.0.1, 10.0.0.2", false},
		{"trailing comma", "10.0.0.1,", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHostList(tt.host)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestBuildConnectionDSN(t *testing.T) {
	dsn := BuildConnectionDSN(Config{
		Host:    "10.0.0.1,10.0.0.2",
		Port:    5432,
		DBName:  "friendlystripe",
		User:    "postgres",
		SSLMode: "disable",
	})
	for _, want := range []string{
		"host=10.0.0.1,10.0.0.2",
		"target_session_attrs=read-write",
		"connect_timeout=4",
	} {
		assert.Contains(t, dsn, want)
	}
}
