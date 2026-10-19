package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func env(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		wantErr string
	}{
		{
			name: "defaults",
			want: Config{Addr: ":8080", LogLevel: slog.LevelInfo, ShutdownTimeout: 10 * time.Second},
		},
		{
			name: "overrides",
			env:  map[string]string{"ADDR": "127.0.0.1:9000", "LOG_LEVEL": "DEBUG", "SHUTDOWN_TIMEOUT": "30s"},
			want: Config{Addr: "127.0.0.1:9000", LogLevel: slog.LevelDebug, ShutdownTimeout: 30 * time.Second},
		},
		{
			name:    "invalid log level",
			env:     map[string]string{"LOG_LEVEL": "loud"},
			wantErr: "invalid LOG_LEVEL",
		},
		{
			name:    "invalid shutdown timeout",
			env:     map[string]string{"SHUTDOWN_TIMEOUT": "soon"},
			wantErr: "invalid SHUTDOWN_TIMEOUT",
		},
		{
			name:    "non positive shutdown timeout",
			env:     map[string]string{"SHUTDOWN_TIMEOUT": "0s"},
			wantErr: "invalid SHUTDOWN_TIMEOUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := Load(env(tt.env))

			// Then
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Load() error = %v, want error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			// And
			if got != tt.want {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
