package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func decryptFile(fp *os.File) ([]byte, error) {
	headerByte := make([]byte, len([]byte(vaultHeaderInfo)))
	_, err := fp.ReadAt(headerByte, 0)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read header: %w", err)
	}

	if !isVaultHeader(headerByte) {
		return nil, ErrorInvalidFormat
	}

	file, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	cypherText, err := parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
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
		return nil, fmt.Errorf("open: %w", err)
	}
	defer fp.Close()

	plainText, err := decryptFile(fp)
	if err != nil {
		return nil, fmt.Errorf("decryptFileAndPrint: %w", err)
	}
	return plainText, nil
}
