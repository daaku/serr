package serr

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

var errFixed = errors.New("fixed")

func TestWrap(t *testing.T) {
	wrappedErr := Wrap(errFixed)
	wrappedErr = Wrap(wrappedErr)
	s := fmt.Sprintf("%+v\n", wrappedErr)
	matched, err := regexp.MatchString(`fixed
github.com.daaku.serr.TestWrap
\s+.*.serr.serr_test.go:14
github.com.daaku.serr.TestWrap
\s+.*.serr.serr_test.go:13
`, s)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("did not match", s)
	}
	we, ok := errors.AsType[*Error](wrappedErr)
	if !ok {
		t.Fatal("expected AsType to work")
	}
	if we.Unwrap() != errFixed {
		t.Fatal("expected errFixed from Unwrap")
	}
	if wrappedErr.Error() != errFixed.Error() {
		t.Fatal("expected Error string to be unchanged")
	}
	if s := fmt.Sprintf("%v", wrappedErr); s != errFixed.Error() {
		t.Fatal("expected v to be errFixed")
	}
	if s := fmt.Sprintf("%s", wrappedErr); s != errFixed.Error() {
		t.Fatal("expected s to be errFixed")
	}
	if s := fmt.Sprintf("%q", wrappedErr); s != `"fixed"` {
		t.Fatal("expected q to be errFixed")
	}
	if len(we.Callers()) < 2 {
		t.Fatal("expected at least 2 callers")
	}
}

func TestErrorf(t *testing.T) {
	err := Errorf("hello")
	s := fmt.Sprintf("%+v\n", err)
	matched, err := regexp.MatchString(`hello
github.com.daaku.serr.TestErrorf
\s+.*.serr.serr_test.go:53
`, s)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("did not match", s)
	}
}

func TestWrapNil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Fatal("nil should be nil")
	}
}
