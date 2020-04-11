package main

import (
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"log"
	"os"
)

var viewCommand = cli.Command{
	Name:   "view",
	Usage:  "View files",
	Flags:  []cli.Flag{},
	Action: viewAction,
}

func viewAction(c *cli.Context) error {
	name := kmsName(
		c.GlobalString("project"),
		c.GlobalString("location"),
		c.GlobalString("keyring"),
		c.GlobalString("key"),
	)

	if len(c.Args()) != 1 {
		log.Fatal("Too many files")
	}
	filename := c.Args().First()

	fstat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	if fstat.IsDir() {
		log.Fatalf("Skipping directory: %s\n", filename)
	}

	err = decryptFileAndPrint(name, filename)
	if err != nil {
		log.Fatal(err)
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
