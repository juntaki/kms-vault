package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
)

const vaultHeaderMagic = "$VAULT;"
const vaultHeaderInfo = vaultHeaderMagic + VaultVersion + ";CLOUD_KMS\n"

var InvalidFormatError = fmt.Errorf("not a vault file")

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
	return bytes.HasPrefix(file, []byte(vaultHeaderMagic))
}
