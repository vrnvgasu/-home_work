package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

func (e Environment) Handle() error {
	for name, envVal := range e {
		if err := os.Unsetenv(name); err != nil {
			return fmt.Errorf("error unsetting environment variable %q: %w", name, err)
		}

		if envVal.NeedRemove {
			continue
		}

		if err := os.Setenv(name, envVal.Value); err != nil {
			return fmt.Errorf("error setting environment variable [%q], value [%s]: %w", name, envVal.Value, err)
		}
	}

	return nil
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("env ReadDir os.ReadDir %s: %w", dir, err)
	}

	envs := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		fileInfo, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("env ReadDir Info: %w", err)
		}
		if fileInfo.Size() == 0 {
			envs[name] = EnvValue{
				NeedRemove: true,
			}

			continue
		}

		value, err := getValue(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("env ReadDir getValue %s: %w", name, err)
		}
		envs[name] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}

	return envs, err
}

func getValue(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("env ReadDir Open %s: %w", filePath, err)
	}
	defer AppClose(file)

	reader := bufio.NewReader(file)
	lineBytes, _, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("env ReadDir ReadLine %s: %w", filePath, err)
	}

	lineBytes = bytes.ReplaceAll(lineBytes, []byte{0}, []byte("\n"))
	line := strings.TrimRight(string(lineBytes), " \v\t")

	return line, nil
}
