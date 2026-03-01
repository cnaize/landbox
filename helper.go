package landbox

import (
	"embed"
	"fmt"
	"os/exec"
)

// NOTE: build from official Linux kernel samples using musl-gcc
// https://github.com/torvalds/linux/tree/v6.19/samples/landlock
//
//go:embed bin/sandboxer
var fsSandboxer embed.FS

func getSandboxer() ([]byte, error) {
	return fsSandboxer.ReadFile("bin/sandboxer")
}

func newCmdError(err error) *exec.Cmd {
	return &exec.Cmd{Err: err}
}

func newCmdErrorf(format string, a ...any) *exec.Cmd {
	return &exec.Cmd{Err: fmt.Errorf(format, a...)}
}
