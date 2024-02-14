package rwatch

import (
	"reflect"
	"testing"
)

func TestNewWatchTree(t *testing.T) {
	type name string
	type input string
	type want struct {
		prefix     string
		globstring string
	}
	tests := []struct {
		name
		input
		want
	}{
		{name("subdir"), input("testdata"), want{"testdata", ""}},
		{name("subdir with trailing slash"), input("testdata/"), want{"testdata/", ""}},
		{name("explicit subdir"), input("./testdata/"), want{"./testdata/", ""}},
		{name("explicit subdir star"), input("./testdata/*"), want{"./testdata", "*"}},
		{name("implicit subdir doublestar"), input("testdata/**"), want{"testdata", "**"}},
		{name("file glob single star"), input("./testdata/*.txt"), want{"./testdata", "*.txt"}},
		{name("file glob double star"), input("./testdata/**.txt"), want{"./testdata", "**.txt"}},
		{name("just movies"), input("testdata/**.avi"), want{"testdata", "**.avi"}},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			gotTree, err := NewWatchTree(string(tt.input))
			if err != nil {
				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if !reflect.DeepEqual(gotTree.prefix, tt.want.prefix) {
				t.Errorf("NewWatchTree() = %v, want %v", gotTree.prefix, tt.want.prefix)
			}
		})
	}
}
