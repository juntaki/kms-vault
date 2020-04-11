package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/xerrors"
)

var vaultHeaderMagic = "$VAULT;"
var vaultHeaderInfo = vaultHeaderMagic + "0.1.0;CLOUD_KMS\n"
var vaultHeaderSize = len([]byte(vaultHeaderInfo))

var InvalidFormatError = xerrors.New("Not a vault file")

func format(ciphertext []byte) []byte {
	return []byte(fmt.Sprintf("%s%s",
		vaultHeaderInfo,
		base64.StdEncoding.EncodeToString(ciphertext)))
}

func parse(file []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(file)
	s := bufio.NewScanner(bytesReader)

	s.Scan() // Discard header
	s.Scan()
	s.Text()
	return base64.StdEncoding.DecodeString(s.Text())
}

func isVaultHeader(file []byte) bool {
	// TODO: Parse version, if file format changes.a
	return bytes.HasPrefix(file, []byte(vaultHeaderMagic))
}
