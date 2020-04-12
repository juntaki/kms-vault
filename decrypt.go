package main

import (
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/xerrors"
)

func decryptFile(fp *os.File) ([]byte, error) {
	headerByte := make([]byte, len([]byte(vaultHeaderInfo)))
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

	plainText, err := kmsClient.Decrypt(cypherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func getPlainText(filename string) ([]byte, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, xerrors.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(fp)
	if err != nil {
		return nil, xerrors.Errorf("decryptFileAndPrint: %w", err)
	}
	return plainText, nil
}
