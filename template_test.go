package main

import (
	"reflect"
	"testing"
)

func Test_convertToTemplateData(t *testing.T) {
	type args struct {
		raw map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				raw: map[string][]byte{"file": []byte("file-contents")},
			},
			want:    map[string]interface{}{"file": "file-contents"},
			wantErr: false,
		},
		{
			name: "yaml",
			args: args{
				raw: map[string][]byte{"file.yaml": []byte("data: file-contents")},
			},
			want:    map[string]interface{}{"file": map[interface{}]interface{}{"data": "file-contents"}},
			wantErr: false,
		},
		{
			name: "duplicate",
			args: args{
				raw: map[string][]byte{
					"file.yaml": []byte("data: file-contents"),
					"file":      []byte("file-contents"),
				},
			},
			wantErr: true,
		},
		{
			name: "merge",
			args: args{
				raw: map[string][]byte{
					"file1.yaml": []byte("data: file-contents"),
					"file2":      []byte("file-contents"),
				},
			},
			want: map[string]interface{}{
				"file1": map[interface{}]interface{}{"data": "file-contents"},
				"file2": "file-contents",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToTemplateData(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToTemplateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToTemplateData() got = %v, want %v", got, tt.want)
			}
		})
	}
}
