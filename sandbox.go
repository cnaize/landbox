package landbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/liamg/memit"
)

const (
	roPathsEnvKey  = "LL_FS_RO"
	rwPathsEnvKey  = "LL_FS_RW"
	tcpListenKey   = "LL_TCP_BIND"
	tcpConnectKey  = "LL_TCP_CONNECT"
	deniedScopeKey = "LL_SCOPED"
)

type Sandbox struct {
	roPaths []string
	rwPaths []string
	sandbox *os.File
	options *Options
}

func NewSandbox(roPaths, rwPaths []string, options *Options) *Sandbox {
	return &Sandbox{
		roPaths: roPaths,
		rwPaths: rwPaths,
		options: options,
	}
}

func (s *Sandbox) Command(name string, arg ...string) *exec.Cmd {
	// lazy init
	if err := s.init(); err != nil {
		return newCmdError(err)
	}

	// prepare command
	cmd := exec.Command(s.sandbox.Name(), append([]string{name}, arg...)...)
	s.prepare(cmd)

	return cmd
}

func (s *Sandbox) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	// lazy init
	if err := s.init(); err != nil {
		return newCmdError(err)
	}

	// prepare command
	cmd := exec.CommandContext(ctx, s.sandbox.Name(), append([]string{name}, arg...)...)
	s.prepare(cmd)

	return cmd
}

func (s *Sandbox) Close() error {
	if s.sandbox != nil {
		return s.sandbox.Close()
	}

	return nil
}

func (s *Sandbox) init() error {
	if s.sandbox != nil {
		return nil
	}

	// get binary file
	sandboxer, err := getSandboxer()
	if err != nil {
		return fmt.Errorf("get sandboxer: %w", err)
	}

	// put file in memory
	_, file, err := memit.Command(bytes.NewReader(sandboxer))
	if err != nil {
		return fmt.Errorf("memit command: %w", err)
	}

	s.sandbox = file

	return nil
}

func (s *Sandbox) prepare(cmd *exec.Cmd) {
	// required
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", roPathsEnvKey, strings.Join(s.roPaths, ":")))
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", rwPathsEnvKey, strings.Join(s.rwPaths, ":")))

	// additional
	if s.options != nil {
		if len(s.options.TCPListen) > 0 {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", tcpListenKey, s.options.TCPListen))
		}
		if len(s.options.TCPConnect) > 0 {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", tcpConnectKey, s.options.TCPConnect))
		}
		if len(s.options.Scope()) > 0 {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", deniedScopeKey, s.options.Scope()))
		}
	}
}
