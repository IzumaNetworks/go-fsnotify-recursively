package fsnotifyr_test

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

func TestNewFolder(t *testing.T) {

	tests := []struct {
		name    string
		fs      fs.FS
		path    string
		parent  fsnotifyr.Folder
		want    string
		wantErr bool
	}{
		{
			"one",
			os.DirFS("testdata"),
			".",
			nil,
			JUST_FOLDERS,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fsnotifyr.NewFolder(tt.fs, tt.path, tt.parent)
			want := strings.TrimSpace(tt.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			os.WriteFile("want.txt", []byte(want), 0644)
			os.WriteFile("got.txt", []byte(got.String()), 0644)

			if !reflect.DeepEqual(got.String(), want) {
				t.Errorf("NewFolder() = %s, want %v", got.String(), want)
			}
		})
	}
}
