package gorph_test

import (
	"io/fs"
	"os"
	"reflect"
	"sync"
	"testing"

	gorph "go.izuma.io/go-fsnotify-recursively"
)

// test that gorph.Root() returns correct root path
func TestNewGorph_Root(t *testing.T) {
	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "self", args: args{root: ".", back: os.DirFS("."), pattern: "*"}, want: ".", wantErr: false},
		{name: "testdata", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "*"}, want: "testdata", wantErr: false},
		{name: "x", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: "testdata", wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: "go.mod", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewWithFs(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got := gotGorph.Root()
				if got != tt.want {
					t.Errorf("NewGorph() = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewGorph() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// test that gorph.Walk() returns entire file tree
func TestNewGorph_Walk(t *testing.T) {
	var VIDEOS = []string{".", "Nuages-3.OGG"}
	var DOWNLOADS = []string{".", "node", "node/node-v20.11.0-linux-x64.tar.xz"}

	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Videos", args: args{root: "testdata/Videos", back: os.DirFS("testdata/Videos"), pattern: "*"}, want: VIDEOS, wantErr: false},
		{name: "Downloads", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "**"}, want: DOWNLOADS, wantErr: false},
		{name: "folder that doesn't exist", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: []string{""}, wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: []string{""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewWithFs(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got, _ := gotGorph.Walk()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewGorph() = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("NewGorph() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// test that glob returns a correctly filtered file tree, with correct usage and single and double star ("*", "**")
func TestNewGorph_Glob(t *testing.T) {
	var VIDEOS = []string{"Nuages-3.OGG"}
	var JUST_NODE = []string{"node"}
	var DOWNLOADS = []string{".", "node", "node/node-v20.11.0-linux-x64.tar.xz"}
	var DOWNLOADS_ROOTED = []string{"Downloads", "Downloads/node", "Downloads/node/node-v20.11.0-linux-x64.tar.xz"}
	var MOMMY = []string{"Documents/the-mommy-book.txt", "Pictures/mommy-and-me.jpeg"}

	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "3 videos", args: args{root: "testdata/Videos", back: os.DirFS("testdata/Videos"), pattern: "*"}, want: VIDEOS, wantErr: false},
		{name: "node", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "*"}, want: JUST_NODE, wantErr: false},
		{name: "Downloads", args: args{root: "testdata/Downloads", back: os.DirFS("testdata/Downloads"), pattern: "**"}, want: DOWNLOADS, wantErr: false},
		{name: "Downloads rooted", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "Downloads/**"}, want: DOWNLOADS_ROOTED, wantErr: false},
		{name: "mommy files", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "**/*mommy*"}, want: MOMMY, wantErr: false},
		{name: "folder that doesn't exist", args: args{root: "x", back: os.DirFS("x"), pattern: "*"}, want: []string{""}, wantErr: true},
		{name: "go.mod", args: args{root: "go.mod", back: os.DirFS("go.mod"), pattern: "*"}, want: []string{""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewWithFs(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got, _ := gotGorph.Glob(tt.args.pattern)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("got %v, wanted %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("got error = %v, but wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// test that gorph.Folders() returns the folder tree
func TestNewGorph_Folders(t *testing.T) {
	var JUST_NODE = []string{".", "Documents", "Documents/torus", "Documents/torus/jamaica", "Downloads", "Downloads/node", "Music", "Pictures", "Videos"}
	type args struct {
		root    string
		back    fs.FS
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "testdata", args: args{root: "testdata", back: os.DirFS("testdata"), pattern: "*"}, want: JUST_NODE, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewWithFs(tt.args.root, tt.args.pattern, tt.args.back)
			if err == nil {
				got := gotGorph.Folders()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("got %v, wanted %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("got error = %v, but wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// test that gorph.Children(parent) returns the children of parent
func TestNewGorph_Children(t *testing.T) {
	var DOCUMENTS = []string{"Documents/torus", "Documents/torus/jamaica"}
	var ROOT = []string{"Documents", "Documents/torus", "Documents/torus/jamaica", "Downloads", "Downloads/node", "Music", "Pictures", "Videos"}

	root := "testdata"
	testData := os.DirFS(root)
	pattern := "*"

	tests := []struct {
		name       string
		childrenOf string
		want       []string
		wantErr    bool
	}{
		{name: "Documents", childrenOf: "Documents", want: DOCUMENTS, wantErr: false},
		{name: "testdata", childrenOf: ".", want: ROOT, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGorph, err := gorph.NewWithFs(root, pattern, testData)
			if err == nil {
				got, _ := gotGorph.Children(tt.childrenOf)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("got %v, wanted %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("got error = %v, but wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNewGorph_Listen(t *testing.T) {

	testData := os.DirFS("testdata")
	expected := "tmp.txt"
	g, err := gorph.NewWithFs("testdata", "**/*.txt", testData)
	if err != nil {
		t.Errorf("error creating new gorph. %v", err)
		return
	}
	evs, _ := g.Listen()
	var ev1 gorph.GorphEvent
	var wg sync.WaitGroup

	//	set up sync
	wg.Add(2)

	//	capture event
	go func() {
		ev1 = <-evs
		wg.Done()
	}()

	//	send event
	go func() {
		os.Create("testdata/" + expected)
		wg.Done()
	}()

	//	wait for events
	wg.Wait()

	//	be nice and clean up
	g.Close()

	//	test the event
	if ev1.Path != "tmp.txt" {
		t.Errorf("expected path of %q but got %q", expected, ev1.Path)
	}

}
