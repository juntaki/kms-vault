package main

import (
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func decryptCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "decrypt",
		Usage:  "Decrypt files",
		Flags:  kmsFlags,
		Action: decryptAction,
	}
}

func decryptAction(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return xerrors.New("Specify at least one file")
	}

	name := kmsNameFromContext(c)
	for _, filename := range c.Args() {
		// Skip dir
		fstat, err := os.Stat(filename)
		if err != nil {
			return err
		}
		if fstat.IsDir() {
			log.Printf("Skipping directory: %s\n", filename)
			continue
		}

		err = decryptFileAndWrite(name, filename)
		if xerrors.Is(err, InvalidFormatError) {
			log.Printf("Skipping not vault file: %s\n", filename)
			return err
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func decryptFileAndWrite(name, filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(name, fp)
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
	log.Printf("Decryption successful: %s\n", filename)
	return nil
}

func decryptFile(name string, fp *os.File) ([]byte, error) {
	headerByte := make([]byte, vaultHeaderSize)
	_, err := fp.ReadAt(headerByte, 0)
	if err != nil && err != io.EOF {
		return nil, xerrors.Errorf("read header: %w", err)
	}

	if !isVaultHeader(headerByte) {
		return nil, InvalidFormatError
	}

	file, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, xerrors.Errorf("readall: %w", err)
	}

	cypherText, err := parse(file)
	if err != nil {
		return nil, xerrors.Errorf("parse: %w", err)
	}

	plainText, err := kmsDecrypt(name, cypherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
