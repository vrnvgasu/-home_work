package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := env.Handle()
	if err != nil {
		log.Printf("Error executing env: %v \n", err)
		return 1
	}

	var arg []string
	switch len(cmd) {
	case 0:
		return 0
	case 1:
		arg = []string{}
	default:
		arg = cmd[1:]
	}

	c := exec.Command(cmd[0], arg...) // #nosec G204
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		log.Printf("Error executing command: %v \n", err)
		return 1
	}

	return c.ProcessState.ExitCode()
}
