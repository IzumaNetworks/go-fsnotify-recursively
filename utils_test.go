package rwatch

import (
	"testing"
)

func TestGetGlobParent(t *testing.T) {
	type name string
	type input string
	type want struct {
		head string
		tail string
	}
	tests := []struct {
		name
		input
		want
	}{
		{name("self"), input("."), want{".", ""}},
		{name("self explicit"), input("./*"), want{".", "*"}},
		{name("self implicit"), input(""), want{"", ""}},
		{name("self explicit doublestar"), input("./**"), want{".", "**"}},
		{name("plain old doublestar"), input("**"), want{"", "**"}},
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
			gotHead, gotTail, err := componentizeGlobString(string(tt.input))
			if err != nil {
				t.Errorf("head:\t%q\ntail:\t%q\nerr:\t%v", gotHead, gotTail, err)
				return
			}
			if gotHead != tt.want.head || gotTail != tt.want.tail {
				t.Errorf("got = %q, %q\t wanted = %q, %q ", gotHead, gotTail, tt.want.head, tt.want.tail)
			}
		})
	}
}
