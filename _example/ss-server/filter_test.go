package main

import "testing"

func TestAccessDenied(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name         string
		args         args
		wantRejected bool
	}{
		{"access denied", args{"www.baidu.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRejected := AccessDenied(tt.args.addr); gotRejected != tt.wantRejected {
				t.Errorf("AccessDenied() = %v, want %v", gotRejected, tt.wantRejected)
			}
		})
	}
}
