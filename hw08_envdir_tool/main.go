package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("not enough arguments, get only [%d]", len(args))
	}

	path, err := filepath.Abs(args[1])
	if err != nil {
		log.Fatal("get abs path:", err)
	}

	var cmd []string
	if len(args) > 2 {
		cmd = args[2:]
	} else {
		cmd = []string{}
	}

	env, err := ReadDir(path)
	if err != nil {
		log.Fatalf("read dir %s: %s", path, err)
	}

	RunCmd(cmd, env)
}
