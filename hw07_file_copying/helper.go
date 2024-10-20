package main

import (
	"fmt"
	"io"
)

func AppClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Printf("Error closing closer: %v", err)
	}
}
