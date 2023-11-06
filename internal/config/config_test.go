package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		osargs  []string
		env     map[string]string
		want    Config
		wantErr bool
	}{
		{
			name:   "Default",
			osargs: []string{"gophermart"},
			env:    map[string]string{},
			want: Config{
				RunAddress:           "localhost:8080",
				AccrualSystemAddress: "http://accrual",
				DarabaseURI:          "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable",
				LogLevel:             "info",
			},
			wantErr: false,
		},
		{
			name:   "Only arguments",
			osargs: []string{"gophermart", "-a", "localhost:8000", "-r", "http://accrual:8000", "-d", "postgres://postgres:postgres@postgres:5432/praktikum", "-l", "debug"},
			env:    map[string]string{},
			want: Config{
				RunAddress:           "localhost:8000",
				AccrualSystemAddress: "http://accrual:8000",
				DarabaseURI:          "postgres://postgres:postgres@postgres:5432/praktikum",
				LogLevel:             "debug",
			},
			wantErr: false,
		},
		{
			name:   "Only environment",
			osargs: []string{"gophermart"},
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8000",
				"ACCRUAL_SYSTEM_ADDRESS": "http://accrual:8000",
				"DATABASE_URI":           "postgres://postgres:postgres@postgres:5432/praktikum",
				"LOG_LEVEL":              "debug",
			},
			want: Config{
				RunAddress:           "localhost:8000",
				AccrualSystemAddress: "http://accrual:8000",
				DarabaseURI:          "postgres://postgres:postgres@postgres:5432/praktikum",
				LogLevel:             "debug",
			},
			wantErr: false,
		},
		{
			name:   "Environment overrides arguments",
			osargs: []string{"gophermart", "-a", "localhost:8000", "-r", "http://accrual:8000", "-d", "postgres://postgres:postgres@postgres:5432/praktikum", "-l", "debug"},
			env: map[string]string{
				"RUN_ADDRESS":            "localhost:8080",
				"ACCRUAL_SYSTEM_ADDRESS": "http://accrual:8080",
				"DATABASE_URI":           "postgres://postgres:postgres@postgres:6432/praktikum",
				"LOG_LEVEL":              "warn",
			},
			want: Config{
				RunAddress:           "localhost:8080",
				AccrualSystemAddress: "http://accrual:8080",
				DarabaseURI:          "postgres://postgres:postgres@postgres:6432/praktikum",
				LogLevel:             "warn",
			},
			wantErr: false,
		},
		{
			name:   "Environment and arguments",
			osargs: []string{"gophermart", "-a", "localhost:8000", "-r", "http://accrual:8000"},
			env: map[string]string{
				"RUN_ADDRESS":  "localhost:8080",
				"DATABASE_URI": "postgres://postgres:postgres@postgres:6432/praktikum",
			},
			want: Config{
				RunAddress:           "localhost:8080",
				AccrualSystemAddress: "http://accrual:8000",
				DarabaseURI:          "postgres://postgres:postgres@postgres:6432/praktikum",
				LogLevel:             "info",
			},
			wantErr: false,
		},
		{
			name:    "With invalid argument",
			osargs:  []string{"gophermart", "-t"},
			env:     map[string]string{},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid argument value",
			osargs:  []string{"gophermart", "-a", "127.0.0.1/8000"},
			env:     map[string]string{},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid evironment variable value",
			osargs:  []string{"gophermart"},
			env:     map[string]string{"RUN_ADDRESS": "127.0.0.1/8000"},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "With invalid evironment variable and argument value",
			osargs:  []string{"gophermart", "-a", "127.0.0.2/8000"},
			env:     map[string]string{"RUN_ADDRESS": "127.0.0.1/8000"},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.osargs
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			c := Config{}
			if err := Parse(&c); err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equal(t, tt.want, c)
		})
	}
}
