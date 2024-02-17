package rwatch

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewFolder(t *testing.T) {

	tests := []struct {
		name    string
		fs      fs.FS
		path    string
		parent  Folder
		want    string
		wantErr bool
	}{
		{
			"one",
			os.DirFS("testdata"),
			".",
			nil, `
.
└── Documents
    ├── mixed
    └── textfiles`,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFolder(tt.fs, tt.path, tt.parent)
			want := strings.TrimSpace(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), want) {
				t.Errorf("NewFolder() = %s, want %v", got.String(), want)
			}
		})
	}
}
