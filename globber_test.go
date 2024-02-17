package fsnotifyr_test

import (
	"reflect"
	"testing"

	fsnotifyr "github.com/sean9999/go-fsnotify-recursively"
)

func TestNewGlobber(t *testing.T) {
	type args struct {
		fullString string
	}
	tests := []struct {
		name    string
		args    args
		want    fsnotifyr.Globber
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fsnotifyr.NewGlobber(tt.args.fullString)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGlobber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGlobber() = %v, want %v", got, tt.want)
			}
		})
	}
}
