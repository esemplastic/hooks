package hooks

import (
	"testing"
)

var messaging = NewHub()

func n() {
	messaging.RunFunc(n)
}

func TestGetCurrentRunner(t *testing.T) {
	messaging.RegisterFunc(n, func() {
		cRunner := GetCurrentRunner().Name
		fnName := ReadSourceFunc(n).Name

		if expected, got := "github.com/esemplastic/hooks.TestGetCurrentRunner", cRunner; expected != got {
			t.Fatalf("expected current runner's name source to be: '%s' but got: '%s'", expected, got)
		}

		if expected, got := "github.com/esemplastic/hooks.n", fnName; expected != got {
			t.Fatalf("expected hook's source to be: '%s' but got: '%s'", expected, got)
		}
	})

	n()
}

var messaging2 = NewHub()

func n2() {
	messaging2.RunFunc(n2)
}

func TestGetCurrentHookSource(t *testing.T) {
	messaging2.RegisterFunc(n2, func() {
		cHookSourceName := GetCurrentHookSource().Name
		fnName := ReadSourceFunc(n2).Name

		if expected, got := "github.com/esemplastic/hooks.n2", cHookSourceName; expected != got {
			t.Fatalf("expected current runner's name source to be: '%s' but got: '%s'", expected, got)
		}

		if expected, got := "github.com/esemplastic/hooks.n2", fnName; expected != got {
			t.Fatalf("expected hook's source to be: '%s' but got: '%s'", expected, got)
		}

	})

	n2()
}
