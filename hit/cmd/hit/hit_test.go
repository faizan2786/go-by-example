package main

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {

	testCases := []struct {
		name    string
		config  *argConfig
		args    []string
		want    *argConfig
		wantErr error
	}{
		{
			name:    "valid",
			config:  &argConfig{},
			args:    []string{"-url=https://www.myserver.com", "-n=5000", "-c=2", "-rps=1000"},
			want:    &argConfig{url: "https://www.myserver.com", n: 5000, c: 2, rps: 1000},
			wantErr: nil,
		},
		{
			name:    "valid_different_order",
			config:  &argConfig{},
			args:    []string{"-n=3000", "-url=https://www.myserver.com", "-rps=1000", "-c=2"},
			want:    &argConfig{url: "https://www.myserver.com", n: 3000, c: 2, rps: 1000},
			wantErr: nil,
		},
		{
			name:    "valid_with_defaults",
			config:  &argConfig{n: 1000, c: 1, rps: 100},
			args:    []string{"-url=https://www.myserver.com"},
			want:    &argConfig{url: "https://www.myserver.com", n: 1000, c: 1, rps: 100},
			wantErr: nil,
		},
		{
			name:    "invalid_argument",
			config:  &argConfig{},
			args:    []string{"-url=https://www.myserver.com", "-n=5000", "-concurrent=2", "-rps=1000"},
			want:    nil,
			wantErr: errors.New("'-concurrent' is not a valid argument"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			gotErr := parseArgs(tt.args, tt.config)

			if tt.wantErr == nil && gotErr != nil {
				t.Fatalf("parseArgs() err = %v, want <nil>", gotErr)
			}

			if tt.wantErr != nil {
				if gotErr == nil {
					t.Fatalf("parseArgs() err = <nil>, want: %v", tt.wantErr)
				}
				if gotErr.Error() != tt.wantErr.Error() {
					t.Fatalf("parseArgs() err = %v, want: %v", gotErr, tt.wantErr)
				}
				return
			}

			if *tt.config != *tt.want {
				t.Fatalf("\ngot:%#v\nwant:%#v\n", tt.config, tt.want)
			}
		})
	}

}
