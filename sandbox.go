package landbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

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
	roPaths Paths
	rwPaths Paths
	options *Options

	once    sync.Once
	sandbox *os.File
}

func NewSandbox(roPaths, rwPaths Paths, options *Options) *Sandbox {
	return &Sandbox{
		roPaths: roPaths,
		rwPaths: rwPaths,
		options: options,
	}
}

func (s *Sandbox) Command(name string, arg ...string) *exec.Cmd {
	// lazy init
	if err := s.init(); err != nil {
		return newCmdErrorf("init failed: %w", err)
	}

	// prepare command
	cmd := exec.Command(s.sandbox.Name(), append([]string{name}, arg...)...)
	s.prepare(cmd)

	return cmd
}

func (s *Sandbox) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	// lazy init
	if err := s.init(); err != nil {
		return newCmdErrorf("init failed: %w", err)
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
	var err error
	s.once.Do(func() {
		// get sandbox file
		sandboxer, err := getSandboxer()
		if err != nil {
			err = fmt.Errorf("get sandboxer: %w", err)
			return
		}

		// load sandbox file
		_, file, err := memit.Command(bytes.NewReader(sandboxer))
		if err != nil {
			err = fmt.Errorf("memit command: %w", err)
			return
		}

		s.sandbox = file
	})

	if s.sandbox == nil {
		err = fmt.Errorf("load sandbox")
	}

	return err
}

func (s *Sandbox) prepare(cmd *exec.Cmd) {
	// required
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", roPathsEnvKey, s.roPaths))
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", rwPathsEnvKey, s.rwPaths))

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
