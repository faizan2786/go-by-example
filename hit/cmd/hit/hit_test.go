package main

import (
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {

	testCases := []struct {
		name       string
		config     *argConfig
		args       []string
		want       *argConfig
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "valid",
			config:  &argConfig{},
			args:    []string{"-n=5000", "-c=2", "-rps=1000", "https://www.myserver.com"},
			want:    &argConfig{url: "https://www.myserver.com", n: 5000, c: 2, rps: 1000},
			wantErr: false,
		},
		{
			name:    "valid_different_order",
			config:  &argConfig{},
			args:    []string{"-n=3000", "-rps=1000", "-c=2", "https://www.myserver.com"},
			want:    &argConfig{url: "https://www.myserver.com", n: 3000, c: 2, rps: 1000},
			wantErr: false,
		},
		{
			name:    "valid_with_defaults",
			config:  &argConfig{c: 1, rps: 100},
			args:    []string{"-n=2000", "https://www.myserver.com"},
			want:    &argConfig{url: "https://www.myserver.com", n: 2000, c: 1, rps: 100},
			wantErr: false,
		},
		{
			name:       "undefined_argument",
			config:     &argConfig{},
			args:       []string{"-n=5000", "-c=2", "-rps=1000", "-url=https://www.myserver.com"},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "flag provided but not defined: -url",
		},
		{
			name:       "invalid_parse_non_int",
			config:     &argConfig{},
			args:       []string{"-n=5000", "-c=two", "-rps=1000", "https://www.myserver.com"},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "invalid value \"two\" for flag -c",
		},
		{
			name:       "invalid_parse_non_positive_int",
			config:     &argConfig{},
			args:       []string{"-n=5000", "-c=0", "-rps=1000", "https://www.myserver.com"},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "invalid value \"0\" for flag -c: value should be greater than 0",
		},
		{
			name:       "invalid_url",
			config:     &argConfig{},
			args:       []string{"-n=5000", "-c=8", "-rps=1000", "www.myserver.com"},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "invalid value \"www.myserver.com\" for url",
		},
		{
			name:       "empty_url",
			config:     &argConfig{},
			args:       []string{"-n=5000", "-c=4", "-rps=1000", ""},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "invalid value \"\" for url",
		},
		{
			name:       "invalid_flag_values",
			config:     &argConfig{},
			args:       []string{"-n=5", "-c=8", "-rps=1000", "https://www.myserver.com"},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "value for flag -c(=8) can not be greater than the value for flag -n(=5)",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			gotErr := parseArgs(tt.args, tt.config)

			if tt.wantErr == false && gotErr != nil {
				t.Fatalf("parseArgs() err = %v, want <nil>", gotErr)
			}

			if tt.wantErr != false {
				if gotErr == nil {
					t.Fatalf("parseArgs() err = <nil>, want: %v", tt.wantErrMsg)
				}
				if !strings.Contains(gotErr.Error(), tt.wantErrMsg) {
					t.Fatalf("parseArgs() err = %v, want: %v", gotErr, tt.wantErrMsg)
				}
				return
			}

			if *tt.config != *tt.want {
				t.Fatalf("\ngot:%#v\nwant:%#v\n", tt.config, tt.want)
			}
		})
	}

}
