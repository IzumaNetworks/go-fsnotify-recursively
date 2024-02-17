package rwatch

import (
	"reflect"
	"strings"
	"testing"
)

const FOLDER_1 string = `.
└── Documents
    ├── mixed
    └── textfiles
	└── torus
	    └── jamaica	
`

func TestNewWatchTree(t *testing.T) {

	type name string
	type input string
	type want string

	tests := []struct {
		name
		input
		want
	}{
		{name("subdir"), input("testdata"), want(strings.TrimSpace(FOLDER_1))},
		{name("subdir with trailing slash"), input("testdata/"), want(strings.TrimSpace(FOLDER_1))},
		{name("explicit subdir"), input("./testdata/"), want(strings.TrimSpace(FOLDER_1))},
		{name("explicit subdir star"), input("./testdata/*"), want(strings.TrimSpace(FOLDER_1))},
		// {name("implicit subdir doublestar"), input("testdata/**"), want{"testdata", "**"}},
		// {name("file glob single star"), input("./testdata/*.txt"), want{"./testdata", "*.txt"}},
		// {name("file glob double star"), input("./testdata/**.txt"), want{"./testdata", "**.txt"}},
		// {name("just movies"), input("testdata/**.avi"), want{"testdata", "**.avi"}},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			gotTree, err := NewWatchTree(string(tt.input))
			gotFolder := gotTree.rootFolder.String()
			if err != nil {

				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if !reflect.DeepEqual(string(tt.want), gotFolder) {
				t.Errorf("wanted %v but got %v", tt.want, gotFolder)
			}
		})
	}
}
