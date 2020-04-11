package main

import (
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func encryptCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "encrypt",
		Usage:  "Encrypt files",
		Flags:  kmsFlags,
		Action: encryptAction,
	}
}

func encryptAction(c *cli.Context) error {
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

		err = encryptFile(name, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func encryptFile(name, filename string) error {
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

	if isVaultHeader(headerByte) {
		log.Printf("Skipping already encrypted: %s\n", filename)
		return nil
	}

	file, err := ioutil.ReadAll(fp)
	if err != nil {
		return xerrors.Errorf("readall: %w", err)
	}

	val, err := kmsEncrypt(
		name,
		file,
	)
	if err != nil {
		return err
	}

	err = fp.Truncate(0)
	if err != nil {
		return xerrors.Errorf("truncate: %w", err)
	}

	_, err = fp.WriteAt(format(val), 0)
	if err != nil {
		return xerrors.Errorf("write: %w", err)
	}
	log.Printf("Encryption successful: %s\n", filename)
	return nil
}
