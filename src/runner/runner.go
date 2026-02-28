package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(args []string) (*Result, error) {
	r, err := run(args)
	if err != nil {
		return r, err
	}
	if r.ExitCode != 0 {
		return r, fmt.Errorf("command %v exited with code %d", args, r.ExitCode)
	}
	return r, nil
}

func RunNoCheck(args []string) (*Result, error) {
	return run(args)
}

func run(args []string) (*Result, error) {
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			return result, nil
		}
		return result, err
	}
	return result, nil
}

// Exec replaces the current process with the given command (like os.execvp).
func Exec(args []string) error {
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, args, os.Environ())
}
