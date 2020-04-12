package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

func checkOverwrite(outputFile string) error {
	_, err := os.Stat(outputFile)
	if err == nil {
		fmt.Print("File exist, overwrite? (y/N): ")
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		text := stdin.Text()
		if !(len(text) > 0 && strings.ToLower(strings.TrimSpace(text))[0] == 'y') {
			return xerrors.New("Aborted")
		}
	}
	return nil
}
