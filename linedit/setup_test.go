package linedit

import (
	"os"
	"testing"
)

func TestIsTerminal(t *testing.T) {
	IsTerminal(os.Stdout.Fd())
	// want := true
	// if got := IsTerminal(os.Stdout.Fd()); got != want {
	// 	t.Errorf("IsTerminal() = %v, want %v", got, want)
	// }
}
