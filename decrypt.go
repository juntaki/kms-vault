package main

import (
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var decryptCommand = cli.Command{
	Name:   "decrypt",
	Usage:  "Decrypt files",
	Flags:  []cli.Flag{},
	Action: decryptAction,
}

func decryptAction(c *cli.Context) error {
	name := kmsName(
		c.GlobalString("project"),
		c.GlobalString("location"),
		c.GlobalString("keyring"),
		c.GlobalString("key"),
	)

	for _, filename := range c.Args() {
		// Skip dir
		fstat, err := os.Stat(filename)
		if err != nil {
			log.Fatal(err)
		}
		if fstat.IsDir() {
			log.Printf("Skipping directory: %s\n", filename)
			continue
		}

		err = decryptFile(name, filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func decryptFile(name, filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
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
		log.Printf("Skipping not vault file: %s\n", filename)
		return nil
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

	err = fp.Truncate(0)
	if err != nil {
		return xerrors.Errorf("truncate: %w", err)
	}

	_, err = fp.WriteAt(plainText, 0)
	if err != nil {
		return xerrors.Errorf("write: %w", err)
	}
	log.Printf("Decryption successful: %s\n", filename)
	return nil
}
