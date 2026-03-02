# Landbox — Landlock "os/exec.Command()" replacement

```go
package main

import "github.com/cnaize/landbox"

func main() {
	// allow only: ro="/usr", rw="/tmp"
	sandbox := landbox.NewSandbox(landbox.Paths{"/usr"}, landbox.Paths{"/tmp"}, nil)
	defer sandbox.Close()

	// deny any other directory
	output, _ := sandbox.Command("ls", "/home").CombinedOutput()

	println(string(output))
	// Executing the sandboxed command...
	// ls: cannot open directory '/home': Permission denied
}
```

# Features:
 - [x] Thread safe
 - [x] Linux amd64 support
 - [ ] Linux arm64 support

# Requirements:
 - Linux kernel 5.13+ (for Landlock LSM support)
