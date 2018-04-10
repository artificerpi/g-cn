package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_transcode(t *testing.T) {
	type args struct {
		r        io.Reader
		encoding bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"transcode_encode_empty", args{strings.NewReader(""), true}, "", false},
		{"transcode_decode_empty", args{strings.NewReader(""), false}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := transcode(tt.args.r, w, tt.args.encoding); (err != nil) != tt.wantErr {
				t.Errorf("transcode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("transcode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
