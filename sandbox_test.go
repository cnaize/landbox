package landbox

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWriteCommand(t *testing.T) {
	var (
		testDirPath  = "testdata"
		testFilePath = filepath.Join(testDirPath, "write.txt")
	)

	tests := []struct {
		name    string
		pass    bool
		roPaths []string
		rwPaths []string
	}{
		{
			name:    "success both paths",
			pass:    true,
			roPaths: []string{"/usr"},
			rwPaths: []string{"testdata"},
		},
		{
			name:    "empty ro paths",
			pass:    false,
			roPaths: nil,
			rwPaths: []string{"testdata"},
		},
		{
			name:    "empty rw paths",
			pass:    false,
			roPaths: []string{"/usr"},
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

			// read with sandbox
			data, err := sandbox.Command("ls", testDirPath).CombinedOutput()
			if tt.pass {
				assert.NoError(t, err, string(data))
			} else {
				assert.Error(t, err, string(data))
			}

			// read without sandbox
			_, err = exec.Command("ls", testDirPath).CombinedOutput()
			assert.NoError(t, err)

			// write with sandbox
			func() {
				defer os.Remove(testFilePath)

				data, err := sandbox.Command("touch", testFilePath).CombinedOutput()
				if tt.pass {
					assert.NoError(t, err, string(data))
				} else {
					assert.Error(t, err, string(data))
				}
			}()

			// write without sandbox
			func() {
				defer os.Remove(testFilePath)

				_, err := exec.Command("touch", testFilePath).CombinedOutput()
				assert.NoError(t, err)
			}()
		})
	}
}
