package main

import (
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"os"
)

func viewCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "view",
		Usage:  "View file",
		Flags:  kmsFlags,
		Action: viewAction,
	}
}

func viewAction(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return xerrors.New("Specify one file")
	}
	filename := c.Args().First()

	name := kmsNameFromContext(c)
	fstat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if fstat.IsDir() {
		return xerrors.Errorf("Skipping directory: %s\n", filename)
	}

	err = decryptFileAndPrint(name, filename)
	if err != nil {
		return err
	}
	return nil
}

func decryptFileAndPrint(name, filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(name, fp)
	if err != nil {
		return xerrors.Errorf("decryptFileAndPrint: %w", err)
	}

	fmt.Println(string(plainText))
	return nil
}
