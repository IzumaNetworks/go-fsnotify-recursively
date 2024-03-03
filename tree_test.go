package fsnotifyr_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

// func Must[RETURNTYPE any](result RETURNTYPE, err error) RETURNTYPE {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return result
// }

// func Must1[ARG any, RETURNTYPE any](fn func(ARG) (RETURNTYPE, error), arg ARG) RETURNTYPE {
// 	r, err := fn(arg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return r
// }

// func Must2[ARG1 any, ARG2 any, RETURNTYPE any](fn func(ARG1, ARG2) (RETURNTYPE, error), arg1 ARG1, arg2 ARG2) RETURNTYPE {
// 	r, err := fn(arg1, arg2)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return r
// }

func Test_folder_FileTree(t *testing.T) {
	type fields struct {
		filesystem fs.FS
		path       string
		parent     fsnotifyr.Folder
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"one", fields{os.DirFS("testdata"), ".", nil}, strings.TrimSpace(EVERYTHING)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folder, err := fsnotifyr.NewFolder(tt.fields.filesystem, tt.fields.path, tt.fields.parent)
			if err != nil {
				t.Error(err)
			}
			got := folder.FileTree(true).String()
			if got != tt.want {
				t.Errorf("got %s but wanted %s", got, tt.want)
			}
		})
	}
}
