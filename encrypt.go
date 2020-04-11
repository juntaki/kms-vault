package main

import (
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"os"
)

var encryptCommand = cli.Command{
	Name:   "encrypt",
	Usage:  "Encrypt files",
	Flags:  []cli.Flag{},
	Action: encryptAction,
}

func encryptAction(c *cli.Context) error {
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

		err = encryptFile(name, filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func encryptFile(name, filename string) error {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()

	file, err := ioutil.ReadAll(fp)
	if isEncrypted(file) {
		log.Printf("Skipping already encrypted: %s\n", filename)
		return nil
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
	_, err = fp.Seek(0, 0)
	if err != nil {
		return xerrors.Errorf("seek: %w", err)
	}

	val = formatEncrypted(val)
	_, err = fp.Write(val)
	if err != nil {
		return xerrors.Errorf("write: %w", err)
	}
	log.Printf("Encryption successful: %s\n", filename)
	return nil
}
