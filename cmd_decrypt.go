package main

import (
	"os"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

func decryptCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "decrypt",
		Usage:  "Decrypt files",
		Flags:  kmsFlags,
		Before: initializeKMS,
		Action: decryptAction,
	}
}

func decryptAction(c *cli.Context) error {
	processed := false
	for _, filename := range c.Args() {
		// Skip dir
		fstat, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if fstat.IsDir() {
			continue
		}

		err = decryptFileAndWrite(filename)
		if err != nil {
			return err
		}
		processed = true
	}

	if !processed {
		return xerrors.New("Specify at least one file")
	}
	return nil
}

func decryptFileAndWrite(filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(fp)
	if err != nil {
		return xerrors.Errorf("decryptFileAndPrint: %w", err)
	}

	err = fp.Truncate(0)
	if err != nil {
		return xerrors.Errorf("truncate: %w", err)
	}

	_, err = fp.WriteAt(plainText, 0)
	if err != nil {
		return xerrors.Errorf("write: %w", err)
	}
	return nil
}
