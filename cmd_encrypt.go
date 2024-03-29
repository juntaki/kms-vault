package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
)

func encryptCommand(kmsFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "encrypt",
		Usage:  "Encrypt files",
		Action: encryptAction,
		Before: initializeKMS,
		Flags:  kmsFlags,
	}
}

func encryptAction(c *cli.Context) error {
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

		err = encryptFileAndWrite(filename)
		if err != nil {
			return err
		}
		processed = true
	}

	if !processed {
		return fmt.Errorf("specify at least one file")
	}
	return nil
}

func encryptFileAndWrite(filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer fp.Close()

	headerByte := make([]byte, len([]byte(vaultHeaderInfo)))
	_, err = fp.ReadAt(headerByte, 0)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read header: %w", err)
	}

	if isVaultHeader(headerByte) {
		log.Printf("Skipping already encrypted: %s\n", filename)
		return nil
	}

	file, err := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("readall: %w", err)
	}

	val, err := kmsClient.Encrypt(file)

	if err != nil {
		return err
	}

	err = fp.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}

	_, err = fp.WriteAt(format(val), 0)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
