package job

import (
	"io"
	"os/exec"
	"sync"
	"time"
)

// Process represents runnable process
type Process struct {
	cmd   *exec.Cmd
	start time.Time
	stop  time.Time
	lock  sync.RWMutex
	stdin io.WriteCloser
	// TODO: Change to whats needed
	state string
}
