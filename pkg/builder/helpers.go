package builder

import (
	"os/exec"
)

func (b *Builder) newCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir    = b.dir
	cmd.Stdout = b.stdout
	cmd.Stderr = b.stderr
	return cmd
}

