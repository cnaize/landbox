package landbox

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteCommands(t *testing.T) {
	var (
		testDirPath  = "testdata"
		testFilePath = filepath.Join(testDirPath, "write.txt")
	)

	tests := []struct {
		name    string
		pass    bool
		roPaths Paths
		rwPaths Paths
	}{
		{
			name:    "success both paths",
			pass:    true,
			roPaths: Paths{"/usr"},
			rwPaths: Paths{"testdata"},
		},
		{
			name:    "empty ro paths",
			pass:    false,
			roPaths: nil,
			rwPaths: Paths{"testdata"},
		},
		{
			name:    "empty rw paths",
			pass:    false,
			roPaths: Paths{"/usr"},
			rwPaths: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create sandbox
			sandbox := NewSandbox(tt.roPaths, tt.rwPaths,
				&Options{
					TCPListen:   Ports{123},
					TCPConnect:  Ports{456, 789},
					DenySockets: true,
					DenySignals: true,
				},
			)
			defer sandbox.Close()

			var wg sync.WaitGroup
			// read with sandbox
			wg.Go(func() {
				data, err := sandbox.Command("ls", testDirPath).CombinedOutput()
				if tt.pass {
					assert.NoError(t, err, string(data))
				} else {
					assert.Error(t, err, string(data))
				}
			})
			// read with sandbox (CommandContext)
			wg.Go(func() {
				data, err := sandbox.CommandContext(context.Background(), "ls", testDirPath).CombinedOutput()
				if tt.pass {
					assert.NoError(t, err, string(data))
				} else {
					assert.Error(t, err, string(data))
				}
			})
			// write with sandbox
			wg.Go(func() {
				defer os.Remove(testFilePath)

				data, err := sandbox.Command("touch", testFilePath).CombinedOutput()
				if tt.pass {
					assert.NoError(t, err, string(data))
				} else {
					assert.Error(t, err, string(data))
				}
			})
			// write with sandbox (CommandContext)
			wg.Go(func() {
				defer os.Remove(testFilePath)

				data, err := sandbox.CommandContext(context.Background(), "touch", testFilePath).CombinedOutput()
				if tt.pass {
					assert.NoError(t, err, string(data))
				} else {
					assert.Error(t, err, string(data))
				}
			})
			// wait till the end
			wg.Wait()

			// read without sandbox
			_, err := exec.Command("ls", testDirPath).CombinedOutput()
			assert.NoError(t, err)

			// write without sandbox
			defer os.Remove(testFilePath)
			_, err = exec.Command("touch", testFilePath).CombinedOutput()
			assert.NoError(t, err)
		})
	}
}
