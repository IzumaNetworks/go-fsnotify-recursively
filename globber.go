package rwatch

import (
	"fmt"
	"os"

	"github.com/gobwas/glob"
)

type Globber interface {
	fmt.Stringer
	glob.Glob
}

type globber struct {
	s string
	g glob.Glob
}

func NewGlobber(s string) Globber {
	g := glob.MustCompile(s, os.PathSeparator)
	return globber{s, g}
}

func (g globber) Match(str string) bool {
	return g.g.Match(str)
}

func (g globber) String() string {
	return g.s
}
