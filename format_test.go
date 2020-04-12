package main

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	val, _ := base64.StdEncoding.DecodeString("CiQAX7XZ1l7f1ImDfuDdJZQ9aFK7a76LlSDtjtGXIbfn53kIK34SLgDKS25uyynLK3OFOMjdPPDtn5dEJtBwUSWOrZwpcGZbwT46rhsrvB574i655Rw=")
	type args struct {
		file []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				file: []byte(`$VAULT;0.1.0;CLOUD_KMS
CiQAX7XZ1l7f1ImDfuDdJZQ9aFK7a76LlSDtjtGXIbfn53kIK34SLgDKS25uyynLK3OFOMjdPPDtn5dEJtBwUSWOrZwpcGZbwT46rhsrvB574i655Rw=`),
			},
			want:    val,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
