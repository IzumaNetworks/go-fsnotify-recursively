package rwatch

import (
	"fmt"
	"os"
	"strings"

	"github.com/gobwas/glob"
)

/**
 *	Globber seperates out the non-magical part of a glob string, compiles a glob, and retains all necessary information
 */

type Globber interface {
	fmt.Stringer
	glob.Glob
}

type globber struct {
	fsRoot   string
	globRoot string
	g        glob.Glob
}

func NewGlobber(fullString string) (Globber, error) {
	fsRoot, globRoot, err := componentizeGlobString(fullString)
	if err != nil {
		return nil, err
	}
	g, err := glob.Compile(globRoot, os.PathSeparator)
	if err != nil {
		return nil, err
	}
	globber := &globber{fsRoot, globRoot, g}
	return globber, nil
}

func (g *globber) Match(str string) bool {
	return g.g.Match(str)
}

func (g *globber) String() string {
	output := map[string]string{
		"fsRoot":   g.fsRoot,
		"globRoot": g.globRoot,
	}
	return fmt.Sprintf("%v", output)
}

func componentizeGlobString(globExpression string) (string, string, error) {
	tail := []string{}
	fullPath := strings.Split(globExpression, string(os.PathSeparator))
	head := fullPath[:]
	for i, slug := range fullPath {
		if isMagic(slug) {
			head = fullPath[:i]
			tail = fullPath[i:]
			break
		}
	}
	return strings.Join(head, string(os.PathSeparator)), strings.Join(tail, string(os.PathSeparator)), nil
}
