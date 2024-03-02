package fsnotifyr_test

import (
	"strings"
	"testing"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

type teatest struct {
	name  string
	input string
	want  string
}

var teas []teatest = []teatest{
	{"subdir", "testdata", NOTHING},
	{"subdir with trailing slash", "testdata/", NOTHING},
	{"explicit subdir", "./testdata/", NOTHING},
	{"explicit subdir star", "./testdata/*", JUST_TOPLEVEL_FOLDERS},
	{"implicit subdir doublestar", "testdata/**", EVERYTHING},
	{"file glob single star", "./testdata/*.txt", NOTHING},
	{"file glob double star", "./testdata/**.txt", ALL_TEXT_FILES},
	{"just movies", "testdata/**.avi", ALL_MOVIES},
}

const JUST_FOLDERS string = `.
├── Documents
│   └── torus
│       └── jamaica
├── Downloads
│   └── node
├── Pictures
└── Videos
`

const JUST_TOPLEVEL_FOLDERS = `.
├── Documents
├── Downloads
├── Pictures
└── Videos
`

const ALL_TEXT_FILES = `.
└── Documents
    ├── alice.txt
    └── pirate.txt
`

const ALL_MOVIES = `.
└── Videos
    ├── VID_1.mp4
    ├── VID_2.mp4
    └── waiting-for-mommy.mov
`

const EVERYTHING = `.
├── Documents
│   ├── Badiou - In Praise of Love.pdf
│   ├── Kristeva - Powers of Horror An Essay on Abjection.pdf
│   ├── alice.txt
│   ├── pirate.txt
│   └── torus
│       └── jamaica
│           └── flag.ico
├── Downloads
│   └── node
│       └── node-v20.11.0-linux-x64.tar.xz
├── Pictures
│   ├── 1-s2.0-S0149763417308692-fx1_lrg.jpg
│   ├── Edward_Hitchcock_Paleontological_Chart.jpeg
│   ├── KruglikovaLikbez.jpeg
│   └── jimenju.png
└── Videos
    ├── VID_1.mp4
    ├── VID_2.mp4
    └── waiting-for-mommy.mov
`

const NOTHING = "."

const FILES_TORUS_DEEP = `.
└── jamaica
    └── flag.ico
`

const FILES_TORUS_SHALLOW = `.
└── jamaica
`

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
		{"just movies", "testdata/**.{avi,mov,mp4}", `{"fsRoot":"testdata","globRoot":"**.{avi,mov,mp4}"}`},
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
		{"subdir", "testdata", NOTHING},
		{"subdir with trailing slash", "testdata/", NOTHING},
		{"explicit subdir", "./testdata/", NOTHING},
		{"explicit subdir star", "./testdata/*", JUST_TOPLEVEL_FOLDERS},
		{"implicit subdir doublestar", "testdata/**", EVERYTHING},
		{"file glob single star", "./testdata/*.txt", NOTHING},
		{"file glob double star", "./testdata/**/*.txt", ALL_TEXT_FILES},
		{"just movies", "testdata/**/*.{mov,mp4,avi}", ALL_MOVIES},
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
