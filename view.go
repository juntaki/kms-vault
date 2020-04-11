package main

import (
	"fmt"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
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

	err = viewFile(name, filename)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func viewFile(name, filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	headerByte := make([]byte, vaultHeaderSize)
	_, err = fp.ReadAt(headerByte, 0)
	if err != nil && err != io.EOF {
		return xerrors.Errorf("read header: %w", err)
	}

	if !isVaultHeader(headerByte) {
		return xerrors.Errorf("not vault file: %s\n", filename)
	}

	file, err := ioutil.ReadAll(fp)
	if err != nil {
		return xerrors.Errorf("readall: %w", err)
	}

	cypherText, err := parse(file)
	if err != nil {
		return xerrors.Errorf("parse: %w", err)
	}

	plainText, err := kmsDecrypt(name, cypherText)
	if err != nil {
		return err
	}

	fmt.Println(string(plainText))
	return nil
}
