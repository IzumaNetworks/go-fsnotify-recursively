package rwatch

import (
	"reflect"
	"testing"
)

func TestNewGlobber(t *testing.T) {
	type args struct {
		fullString string
	}
	tests := []struct {
		name    string
		args    args
		want    Globber
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGlobber(tt.args.fullString)
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
