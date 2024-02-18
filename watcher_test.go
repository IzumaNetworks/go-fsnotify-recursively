package fsnotifyr_test

import (
	"strings"
	"testing"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

const FOLDER_1 string = `.
└── Documents
    ├── mixed
    ├── textfiles
    └── torus
        └── jamaica
`

const FILES_JUST_DOCUMENTS = `.
└── Documents
`

const FILES_ALL_TEXT = `.
└── Documents
    ├── mixed
    │   └── pho.txt
    ├── narf.txt
    ├── textfiles
    │   ├── fi.txt
    │   └── foo.txt
    └── torus
        └── jamaica
            └── foo.txt
`

const FILES_ALL_MOVIES = `.
├── Documents
│   └── mixed
│       └── fum.avi
└── blarg.avi
`

const FILES_ALL = `.
├── Documents
│   ├── mixed
│   │   ├── fum.avi
│   │   └── pho.txt
│   ├── narf.txt
│   ├── textfiles
│   │   ├── fi.txt
│   │   └── foo.txt
│   └── torus
│       └── jamaica
│           └── foo.txt
└── blarg.avi
`

const FILES_NONE = "."

const FILES_TORUS_DEEP = `.
└── jamaica
    └── foo.txt
`

const FILES_TORUS_SHALLOW = `.
└── jamaica
`

func TestNewWatchTree(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"subdir", "testdata", FOLDER_1},
		{"subdir with trailing slash", "testdata/", FOLDER_1},
		{"explicit subdir", "./testdata/", FOLDER_1},
		{"explicit subdir star", "./testdata/*", FOLDER_1},
		{"implicit subdir doublestar", "testdata/**", FOLDER_1},
		{"file glob single star", "./testdata/*.txt", FOLDER_1},
		{"file glob double star", "./testdata/**.txt", FOLDER_1},
		{"just movies", "testdata/**.avi", FOLDER_1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchTree, err := fsnotifyr.NewWatchTree(string(tt.input))
			gotFolder := watchTree.RootFolder().String()
			if err != nil {
				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if gotFolder != strings.TrimSpace(tt.want) {
				t.Errorf("wanted %v but got %v", tt.want, gotFolder)
			}
		})
	}
}

func TestNewWatchTree2(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"subdir", "testdata", FILES_NONE},
		{"subdir with trailing slash", "testdata/", FILES_NONE},
		{"explicit subdir", "./testdata/", FILES_NONE},
		{"explicit subdir star", "./testdata/*", FILES_JUST_DOCUMENTS},
		{"implicit subdir doublestar", "testdata/**", FILES_ALL},
		{"file glob single star", "./testdata/*.txt", FILES_NONE},
		{"file glob double star", "./testdata/**.txt", FILES_ALL_TEXT},
		{"just movies", "testdata/**.avi", FILES_ALL_MOVIES},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchTree, err := fsnotifyr.NewWatchTree(string(tt.input))
			gotFolder := watchTree.RootFolder().FileTree(true).String()
			if err != nil {
				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if gotFolder != strings.TrimSpace(tt.want) {
				t.Errorf("wanted %v but got %v", tt.want, gotFolder)
			}
		})
	}
}

func TestNewWatchTree_Globber(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"subdir", "testdata", `{"fsRoot":"testdata","globRoot":""}`},
		{"subdir with trailing slash", "testdata/", `{"fsRoot":"testdata/","globRoot":""}`},
		{"explicit subdir", "./testdata/", `{"fsRoot":"./testdata/","globRoot":""}`},
		{"explicit subdir star", "./testdata/*", `{"fsRoot":"./testdata","globRoot":"*"}`},
		{"implicit subdir doublestar", "testdata/**", `{"fsRoot":"testdata","globRoot":"**"}`},
		{"file glob single star", "./testdata/*.txt", `{"fsRoot":"./testdata","globRoot":"*.txt"}`},
		{"file glob double star", "./testdata/**.txt", `{"fsRoot":"./testdata","globRoot":"**.txt"}`},
		{"just movies", "testdata/**.avi", `{"fsRoot":"testdata","globRoot":"**.avi"}`},
		{"torus shallow", "testdata/Documents/torus/*", `{"fsRoot":"testdata/Documents/torus","globRoot":"*"}`},
		{"torus deep", "testdata/Documents/torus/**", `{"fsRoot":"testdata/Documents/torus","globRoot":"**"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchTree, err := fsnotifyr.NewWatchTree(string(tt.input))
			gotGlob := watchTree.Globber().String()
			if err != nil {
				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if gotGlob != strings.TrimSpace(tt.want) {
				t.Errorf("wanted %v but got %v", tt.want, gotGlob)
			}
		})
	}
}

func TestNewWatchTree_GlobTree(t *testing.T) {

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"subdir", "testdata", FILES_NONE},
		{"subdir with trailing slash", "testdata/", FILES_NONE},
		{"explicit subdir", "./testdata/", FILES_NONE},
		{"explicit subdir star", "./testdata/*", FILES_JUST_DOCUMENTS},
		{"implicit subdir doublestar", "testdata/**", FILES_ALL},
		{"file glob single star", "./testdata/*.txt", FILES_NONE},
		{"file glob double star", "./testdata/**/*.txt", FILES_ALL_TEXT},
		{"just movies", "testdata/**/*.avi", FILES_ALL_MOVIES},
		{"torus shallow", "testdata/Documents/torus/*", FILES_TORUS_SHALLOW},
		{"torus deep", "testdata/Documents/torus/**", FILES_TORUS_DEEP},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchTree, err := fsnotifyr.NewWatchTree(string(tt.input))
			gotGlobTree := watchTree.RootFolder().GlobTree(watchTree.Globber()).String()
			if err != nil {
				t.Errorf("NewWatchTree() error = %v, input %v", err, tt.input)
				return
			}
			if gotGlobTree != strings.TrimSpace(tt.want) {
				t.Errorf("wanted %v but got %v", tt.want, gotGlobTree)
			}
		})
	}
}
